package icon

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
	"gopkg.in/yaml.v3"
)

var (
	ErrNoIcon         = errors.New("icon entry must have icon set")
	ErrInvalidRegex   = errors.New("invalid regex pattern")
	ErrConfigNotFound = errors.New("user config not found")
)

type ConfigSection struct {
	FallbackIcon string
	ShowName     bool
}

type IconEntry struct {
	DisplayName string
	Icon        string
	Pattern     *regexp.Regexp
	Key         string
}

type Config struct {
	Config ConfigSection
	Icons  map[string]IconEntry
}

type rawConfig struct {
	Config rawConfigSection `yaml:"config"`
	Icons  map[string]any   `yaml:"icons"`
}

type rawConfigSection struct {
	FallbackIcon string `yaml:"fallback-icon"`
	ShowName     bool   `yaml:"show-name"`
}

// buildPattern creates a regex pattern from a key.
// The pattern matches from the start of the command string.
// This allows "nvim" to match "nvim", "nvim -c DBUI", "nvim foo", etc.
// More specific patterns (longer matches) will win via longest-match selection.
func buildPattern(key string) string {
	return "^" + regexp.QuoteMeta(key)
}

// parseEntry parses a YAML value into an IconEntry.
// Supports both string values (icon only) and object values (name and/or icon).
func parseEntry(key string, raw any) (IconEntry, error) {
	entry := IconEntry{Key: key}

	// Handle string value: icon only
	if iconStr, ok := raw.(string); ok {
		entry.Icon = iconStr
		entry.DisplayName = key
		pattern := buildPattern(key)
		re, err := regexp.Compile(pattern)
		if err != nil {
			return IconEntry{}, err
		}
		entry.Pattern = re
		return entry, nil
	}

	// Handle object value (map[string]interface{} from yaml.Unmarshal)
	objMap, ok := raw.(map[string]interface{})
	if !ok {
		return IconEntry{}, errors.New("invalid icon entry format")
	}

	name := ""
	icon := ""
	if n, ok := objMap["name"].(string); ok {
		name = n
	}
	if i, ok := objMap["icon"].(string); ok {
		icon = i
	}

	return parseObjectEntry(key, name, icon)
}

func parseObjectEntry(key, name, icon string) (IconEntry, error) {
	entry := IconEntry{Key: key}

	if icon == "" {
		return IconEntry{}, ErrNoIcon
	}
	entry.Icon = icon

	if name != "" {
		entry.DisplayName = name
	} else {
		entry.DisplayName = key
	}

	pattern := buildPattern(key)
	re, err := regexp.Compile(pattern)
	if err != nil {
		return IconEntry{}, err
	}
	entry.Pattern = re

	return entry, nil
}

//go:embed icons.yaml
var embeddedConfig []byte

func LoadConfig() (*Config, error) {
	embedded, err := parseAndValidate(embeddedConfig)
	if err != nil {
		return nil, err
	}

	userConfig, err := loadUserConfig()
	if err != nil && !errors.Is(err, ErrConfigNotFound) {
		return nil, err
	}

	if userConfig == nil {
		return embedded, nil
	}

	return merge(userConfig, embedded), nil
}

func loadUserConfig() (*Config, error) {
	configHome := env.XdgConfigHome
	if configHome == "" {
		configHome = filepath.Join(env.Home, ".config")
	}

	userPath := filepath.Join(configHome, "tmux", "icons.yaml")
	data, err := os.ReadFile(userPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}

	return parseAndValidate(data)
}

func parseAndValidate(data []byte) (*Config, error) {
	var raw rawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	cfg := &Config{
		Config: ConfigSection{
			FallbackIcon: raw.Config.FallbackIcon,
			ShowName:     raw.Config.ShowName,
		},
		Icons: make(map[string]IconEntry),
	}

	for key, rawEntry := range raw.Icons {
		entry, err := parseEntry(key, rawEntry)
		if err != nil {
			return nil, fmt.Errorf("parsing entry %q: %w", key, err)
		}
		cfg.Icons[key] = entry
	}

	return cfg, nil
}

func merge(user, embedded *Config) *Config {
	cfg := &Config{
		Config: embedded.Config,
		Icons:  make(map[string]IconEntry),
	}

	for k, v := range embedded.Icons {
		cfg.Icons[k] = v
	}

	if user.Config.FallbackIcon != "" {
		cfg.Config.FallbackIcon = user.Config.FallbackIcon
	}
	cfg.Config.ShowName = user.Config.ShowName || embedded.Config.ShowName

	for k, v := range user.Icons {
		cfg.Icons[k] = v
	}

	return cfg
}

type matchResult struct {
	entry       IconEntry
	matchLength int
	keyLen      int
}

func Match(cfg *Config, currentCommand, normalizedCmd string) IconEntry {
	matches := matchAll(cfg, currentCommand, normalizedCmd)
	if len(matches) == 0 {
		return IconEntry{}
	}
	return selectWinner(matches)
}

func matchAll(cfg *Config, currentCommand, normalizedCmd string) []matchResult {
	var matches []matchResult

	for _, entry := range cfg.Icons {
		loc := entry.Pattern.FindStringIndex(normalizedCmd)
		if loc != nil {
			matches = append(matches, matchResult{
				entry:       entry,
				matchLength: loc[1] - loc[0],
				keyLen:      len(entry.Key),
			})
		}
	}

	return matches
}

func selectWinner(matches []matchResult) IconEntry {
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].matchLength != matches[j].matchLength {
			return matches[i].matchLength > matches[j].matchLength
		}
		if matches[i].keyLen != matches[j].keyLen {
			return matches[i].keyLen > matches[j].keyLen
		}
		return strings.Compare(matches[i].entry.Key, matches[j].entry.Key) < 0
	})

	return matches[0].entry
}
