package icon

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
	"gopkg.in/yaml.v3"
)

var (
	ErrNoIconOrPattern = errors.New("icon entry must have at least one of icon or pattern set")
	ErrInvalidPattern  = errors.New("invalid regex pattern")
	ErrConfigNotFound  = errors.New("user config not found")
)

type ConfigSection struct {
	FallbackIcon string
	ShowName     bool
}

type IconEntry struct {
	DisplayName string
	Icon        string
	Pattern     *regexp.Regexp
}

type Config struct {
	Config ConfigSection
	Icons  map[string]IconEntry
}

type rawConfig struct {
	Config rawConfigSection        `yaml:"config"`
	Icons  map[string]rawIconEntry `yaml:"icons"`
}

type rawConfigSection struct {
	FallbackIcon string `yaml:"fallback-icon"`
	ShowName     bool   `yaml:"show-name"`
}

type rawIconEntry struct {
	Icon    string `yaml:"icon"`
	Pattern string `yaml:"pattern"`
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
		entry := IconEntry{DisplayName: key, Icon: rawEntry.Icon}

		if rawEntry.Pattern != "" {
			re, err := regexp.Compile(rawEntry.Pattern)
			if err != nil {
				return nil, err
			}
			entry.Pattern = re
		}

		if err := validateEntry(entry); err != nil {
			return nil, err
		}

		cfg.Icons[key] = entry
	}

	return cfg, nil
}

func validateEntry(entry IconEntry) error {
	if entry.Icon == "" && entry.Pattern == nil {
		return ErrNoIconOrPattern
	}
	return nil
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

	basename := currentCommand
	if i := strings.Index(basename, " "); i >= 0 {
		basename = basename[:i]
	}

	for _, entry := range cfg.Icons {
		if entry.Pattern != nil {
			loc := entry.Pattern.FindStringIndex(normalizedCmd)
			if loc != nil {
				matches = append(matches, matchResult{
					entry:       entry,
					matchLength: loc[1] - loc[0],
					keyLen:      len(entry.DisplayName),
				})
			}
		} else {
			if entry.DisplayName == basename || entry.DisplayName == currentCommand {
				matches = append(matches, matchResult{
					entry:       entry,
					matchLength: 0,
					keyLen:      len(entry.DisplayName),
				})
			}
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
		return strings.Compare(matches[i].entry.DisplayName, matches[j].entry.DisplayName) < 0
	})

	return matches[0].entry
}
