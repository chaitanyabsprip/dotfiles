package icon

import (
	"regexp"
	"testing"
)

func TestBuildPattern(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		{"nvim", "^nvim"},
		{"lazygit", "^lazygit"},
		{"nvim -c DBUI", "^nvim -c DBUI"},
		{"lazygit|gitui", "^lazygit\\|gitui"},
		{"?", "^\\?"},
	}

	for _, tt := range tests {
		got := buildPattern(tt.key)
		if got != tt.expected {
			t.Errorf("buildPattern(%q) = %q; want %q", tt.key, got, tt.expected)
		}
	}
}

func TestMatch_ExactMatch(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim": {Key: "nvim", DisplayName: "nvim", Icon: "N", Pattern: mustCompile("^nvim")},
		},
	}

	got := Match(cfg, "nvim", "nvim")
	if got.DisplayName != "nvim" || got.Icon != "N" {
		t.Errorf("Match(nvim) = %+v; want {DisplayName: nvim, Icon: N}", got)
	}
}

func TestMatch_ExactMatchWithArgs(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			// Pattern matches "nvim" at the start of "nvim -c DBUI"
			"nvim": {Key: "nvim", DisplayName: "nvim", Icon: "N", Pattern: mustCompile("^nvim")},
		},
	}

	got := Match(cfg, "nvim", "nvim -c DBUI")
	if got.DisplayName != "nvim" || got.Icon != "N" {
		t.Errorf("Match(nvim with args) = %+v; want {DisplayName: nvim, Icon: N}", got)
	}
}

func TestMatch_ComplexKey(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim -c DBUI": {Key: "nvim -c DBUI", DisplayName: "db", Icon: "D", Pattern: mustCompile("nvim -c DBUI")},
		},
	}

	got := Match(cfg, "nvim -c DBUI", "nvim -c DBUI")
	if got.DisplayName != "db" || got.Icon != "D" {
		t.Errorf("Match(nvim -c DBUI) = %+v; want {DisplayName: db, Icon: D}", got)
	}
}

func TestMatch_LongestMatchWins(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim":         {Key: "nvim", DisplayName: "nvim", Icon: "N", Pattern: mustCompile("^nvim")},
			"nvim -c DBUI": {Key: "nvim -c DBUI", DisplayName: "db", Icon: "D", Pattern: mustCompile("nvim -c DBUI")},
		},
	}

	got := Match(cfg, "nvim -c DBUI", "nvim -c DBUI")
	if got.DisplayName != "db" {
		t.Errorf("Match(nvim -c DBUI) = %+v; want {DisplayName: db} (longer match)", got)
	}
}

func TestMatch_SameNameDifferentIcons(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"claude":   {Key: "claude", DisplayName: "agent", Icon: "C", Pattern: mustCompile("^claude$")},
			"opencode": {Key: "opencode", DisplayName: "agent", Icon: "O", Pattern: mustCompile("^opencode$")},
		},
	}

	gotClaude := Match(cfg, "claude", "claude")
	if gotClaude.DisplayName != "agent" || gotClaude.Icon != "C" {
		t.Errorf("Match(claude) = %+v; want {DisplayName: agent, Icon: C}", gotClaude)
	}

	gotOpencode := Match(cfg, "opencode", "opencode")
	if gotOpencode.DisplayName != "agent" || gotOpencode.Icon != "O" {
		t.Errorf("Match(opencode) = %+v; want {DisplayName: agent, Icon: O}", gotOpencode)
	}
}

func TestMatch_NoMatchReturnsZeroValue(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim": {Key: "nvim", DisplayName: "nvim", Icon: "N", Pattern: mustCompile("^nvim")},
		},
	}

	got := Match(cfg, "unknown", "unknown")
	if got != (IconEntry{}) {
		t.Errorf("Match(unknown) = %+v; want zero-value IconEntry", got)
	}
}

func TestMatch_TieBreakByKeyLength(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim":   {Key: "nvim", DisplayName: "short", Icon: "S", Pattern: mustCompile("nvim")},
			"neovim": {Key: "neovim", DisplayName: "long", Icon: "L", Pattern: mustCompile("neovim")},
		},
	}

	got := Match(cfg, "neovim", "neovim")
	if got.DisplayName != "long" {
		t.Errorf("Match(neovim) = %+v; want {DisplayName: long} (longer key)", got)
	}
}

func TestMatch_TieBreakAlphabetical(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"zebra": {Key: "zebra", DisplayName: "z", Icon: "1", Pattern: mustCompile("zebra")},
			"alpha": {Key: "alpha", DisplayName: "a", Icon: "2", Pattern: mustCompile("alpha")},
		},
	}

	got := Match(cfg, "alpha", "alpha")
	if got.DisplayName != "a" {
		t.Errorf("Match(alpha) = %+v; want {DisplayName: a} (alpha wins tie)", got)
	}
}

func TestParseEntry_StringValue(t *testing.T) {
	entry, err := parseEntry("nvim", "N")
	if err != nil {
		t.Fatalf("parseEntry(nvim, N) error = %v", err)
	}
	if entry.Icon != "N" {
		t.Errorf("Icon = %q; want N", entry.Icon)
	}
	if entry.DisplayName != "nvim" {
		t.Errorf("DisplayName = %q; want nvim", entry.DisplayName)
	}
	if entry.Key != "nvim" {
		t.Errorf("Key = %q; want nvim", entry.Key)
	}
	if entry.Pattern.String() != "^nvim" {
		t.Errorf("Pattern = %v; want ^nvim", entry.Pattern)
	}
}

func TestParseEntry_ObjectWithName(t *testing.T) {
	entry, err := parseEntry("lazygit", map[string]any{"name": "git", "icon": "G"})
	if err != nil {
		t.Fatalf("parseEntry(lazygit, ...) error = %v", err)
	}
	if entry.Icon != "G" {
		t.Errorf("Icon = %q; want G", entry.Icon)
	}
	if entry.DisplayName != "git" {
		t.Errorf("DisplayName = %q; want git", entry.DisplayName)
	}
	if entry.Key != "lazygit" {
		t.Errorf("Key = %q; want lazygit", entry.Key)
	}
}

func TestParseEntry_ObjectWithoutName(t *testing.T) {
	entry, err := parseEntry("nvim -c DBUI", map[string]any{"icon": "D", "name": "db"})
	if err != nil {
		t.Fatalf("parseEntry(nvim -c DBUI, ...) error = %v", err)
	}
	if entry.DisplayName != "db" {
		t.Errorf("DisplayName = %q; want db", entry.DisplayName)
	}
}

func TestParseEntry_MissingIcon(t *testing.T) {
	_, err := parseEntry("nvim", map[string]any{"name": "editor"})
	if err != ErrNoIcon {
		t.Errorf("parseEntry error = %v; want ErrNoIcon", err)
	}
}

func TestParseAndValidate_SimpleEntry(t *testing.T) {
	data := []byte(`
icons:
  nvim: "N"
`)
	cfg, err := parseAndValidate(data)
	if err != nil {
		t.Fatalf("parseAndValidate error = %v", err)
	}
	entry, ok := cfg.Icons["nvim"]
	if !ok {
		t.Fatal("nvim entry not found")
	}
	if entry.Icon != "N" {
		t.Errorf("Icon = %q; want N", entry.Icon)
	}
	if entry.DisplayName != "nvim" {
		t.Errorf("DisplayName = %q; want nvim", entry.DisplayName)
	}
	if entry.Pattern.String() != "^nvim" {
		t.Errorf("Pattern = %v; want ^nvim", entry.Pattern)
	}
}

func TestParseAndValidate_ComplexEntry(t *testing.T) {
	data := []byte(`
icons:
  "nvim -c DBUI":
    name: "db"
    icon: "D"
`)
	cfg, err := parseAndValidate(data)
	if err != nil {
		t.Fatalf("parseAndValidate error = %v", err)
	}
	entry, ok := cfg.Icons["nvim -c DBUI"]
	if !ok {
		t.Fatal("nvim -c DBUI entry not found")
	}
	if entry.Icon != "D" {
		t.Errorf("Icon = %q; want D", entry.Icon)
	}
	if entry.DisplayName != "db" {
		t.Errorf("DisplayName = %q; want db", entry.DisplayName)
	}
}

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error = %v", err)
	}

	if cfg.Config.FallbackIcon != "?" {
		t.Errorf("FallbackIcon = %q; want ?", cfg.Config.FallbackIcon)
	}

	entry, ok := cfg.Icons["nvim"]
	if !ok {
		t.Fatal("nvim entry not found")
	}
	if entry.Icon == "" {
		t.Error("nvim Icon is empty")
	}

	t.Logf("Total icons: %d", len(cfg.Icons))
}

func TestMerge_UserOverridesEmbedded(t *testing.T) {
	embedded := &Config{
		Config: ConfigSection{FallbackIcon: "?"},
		Icons: map[string]IconEntry{
			"nvim": {Key: "nvim", DisplayName: "nvim", Icon: "N", Pattern: mustCompile("^nvim")},
		},
	}
	user := &Config{
		Config: ConfigSection{FallbackIcon: "X"},
		Icons: map[string]IconEntry{
			"nvim": {Key: "nvim", DisplayName: "nvim", Icon: "X", Pattern: mustCompile("^nvim")},
		},
	}

	got := merge(user, embedded)

	if got.Config.FallbackIcon != "X" {
		t.Errorf("FallbackIcon = %q; want X", got.Config.FallbackIcon)
	}
	if got.Icons["nvim"].Icon != "X" {
		t.Errorf("nvim Icon = %q; want X", got.Icons["nvim"].Icon)
	}
}

func TestMerge_EmbeddedEntriesPreserved(t *testing.T) {
	embedded := &Config{
		Config: ConfigSection{FallbackIcon: "?"},
		Icons: map[string]IconEntry{
			"nvim": {Key: "nvim", DisplayName: "nvim", Icon: "N", Pattern: mustCompile("^nvim")},
			"git":  {Key: "git", DisplayName: "git", Icon: "G", Pattern: mustCompile("^git$")},
		},
	}
	user := &Config{
		Config: ConfigSection{FallbackIcon: "X"},
		Icons: map[string]IconEntry{
			"nvim": {Key: "nvim", DisplayName: "nvim", Icon: "X", Pattern: mustCompile("^nvim")},
		},
	}

	got := merge(user, embedded)

	if _, ok := got.Icons["git"]; !ok {
		t.Error("git entry not preserved")
	}
	if got.Icons["nvim"].Icon != "X" {
		t.Error("nvim not overridden")
	}
}

func mustCompile(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}
