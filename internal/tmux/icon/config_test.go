package icon

import (
	"regexp"
	"testing"
)

func TestMatch_ExactKeyMatch(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim": {DisplayName: "nvim", Icon: "N"},
		},
	}

	got := Match(cfg, "nvim", "nvim")
	if got.DisplayName != "nvim" || got.Icon != "N" {
		t.Errorf("Match(nvim) = %+v; want {DisplayName: nvim, Icon: N}", got)
	}
}

func TestMatch_ExactKeyMatchWithArgs(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim": {DisplayName: "nvim", Icon: "N"},
		},
	}

	got := Match(cfg, "nvim", "nvim mongo_to_postgres.py")
	if got.DisplayName != "nvim" || got.Icon != "N" {
		t.Errorf("Match(nvim with args) = %+v; want {DisplayName: nvim, Icon: N}", got)
	}
}

func TestMatch_PatternMatch(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"db": {
				DisplayName: "db",
				Icon:        "D",
				Pattern:     regexp.MustCompile(`nvim -c DBUI`),
			},
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
			"short": {
				DisplayName: "short",
				Icon:        "S",
				Pattern:     regexp.MustCompile(`nvim`),
			},
			"long": {
				DisplayName: "long",
				Icon:        "L",
				Pattern:     regexp.MustCompile(`nvim -c DBUI`),
			},
		},
	}

	got := Match(cfg, "nvim -c DBUI", "nvim -c DBUI")
	if got.DisplayName != "long" {
		t.Errorf("Match(nvim -c DBUI) = %+v; want {DisplayName: long} (longest match)", got)
	}
}

func TestMatch_NoMatchReturnsZeroValue(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim": {DisplayName: "nvim", Icon: "N"},
		},
	}

	got := Match(cfg, "unknown-cmd", "unknown-cmd")
	if got != (IconEntry{}) {
		t.Errorf("Match(unknown-cmd) = %+v; want zero-value IconEntry", got)
	}
}

func TestMatch_NoPatternLosesToPattern(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim": {DisplayName: "nvim", Icon: "N"},
			"db": {
				DisplayName: "db",
				Icon:        "D",
				Pattern:     regexp.MustCompile(`nvim`),
			},
		},
	}

	got := Match(cfg, "nvim -c DBUI", "nvim -c DBUI")
	if got.DisplayName != "db" {
		t.Errorf("Match(nvim -c DBUI) = %+v; want {DisplayName: db} (pattern beats no-pattern)", got)
	}
}

func TestMatch_TieBreakByKeyLength(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"a": {
				DisplayName: "a",
				Icon:        "1",
				Pattern:     regexp.MustCompile(`nvim`),
			},
			"abc": {
				DisplayName: "abc",
				Icon:        "2",
				Pattern:     regexp.MustCompile(`nvim`),
			},
		},
	}

	got := Match(cfg, "nvim -c DBUI", "nvim -c DBUI")
	if got.DisplayName != "abc" {
		t.Errorf("Match(nvim -c DBUI) = %+v; want {DisplayName: abc} (longer key wins tie)", got)
	}
}

func TestMatch_TieBreakAlphabetical(t *testing.T) {
	cfg := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"zebra": {
				DisplayName: "zebra",
				Icon:        "1",
				Pattern:     regexp.MustCompile(`nvim`),
			},
			"alpha": {
				DisplayName: "alpha",
				Icon:        "2",
				Pattern:     regexp.MustCompile(`nvim`),
			},
		},
	}

	got := Match(cfg, "nvim -c DBUI", "nvim -c DBUI")
	if got.DisplayName != "alpha" {
		t.Errorf("Match(nvim -c DBUI) = %+v; want {DisplayName: alpha} (alpha wins tie)", got)
	}
}

func TestNormalizeCommand(t *testing.T) {
	tests := []struct {
		raw  string
		want string
	}{
		{"/usr/local/bin/nvim -c DBUI", "nvim -c DBUI"},
		{"/opt/homebrew/bin/nvim", "nvim"},
		{"nvim", "nvim"},
		{"nvim -c DBUI", "nvim -c DBUI"},
		{"", ""},
		{"vim", "vim"},
	}

	for _, tt := range tests {
		got := normalizeCommand(tt.raw)
		if got != tt.want {
			t.Errorf("normalizeCommand(%q) = %q; want %q", tt.raw, got, tt.want)
		}
	}
}

func TestFormat_ShowNameTrueWithIcon(t *testing.T) {
	cfg := ConfigSection{FallbackIcon: "?", ShowName: true}
	entry := IconEntry{DisplayName: "nvim", Icon: "N"}

	got := Format(entry, cfg, "fallback")
	want := "N nvim"
	if got != want {
		t.Errorf("Format() = %q; want %q", got, want)
	}
}

func TestFormat_ShowNameTrueWithoutIcon(t *testing.T) {
	cfg := ConfigSection{FallbackIcon: "?", ShowName: true}
	entry := IconEntry{DisplayName: "nvim", Icon: ""}

	got := Format(entry, cfg, "fallback")
	want := "nvim"
	if got != want {
		t.Errorf("Format() = %q; want %q", got, want)
	}
}

func TestFormat_ShowNameFalse(t *testing.T) {
	cfg := ConfigSection{FallbackIcon: "?", ShowName: false}
	entry := IconEntry{DisplayName: "nvim", Icon: "N"}

	got := Format(entry, cfg, "fallback")
	want := "N"
	if got != want {
		t.Errorf("Format() = %q; want %q", got, want)
	}
}

func TestFormat_Fallback(t *testing.T) {
	cfg := ConfigSection{FallbackIcon: "?", ShowName: true}
	entry := IconEntry{}

	got := Format(entry, cfg, "zsh")
	want := "? zsh"
	if got != want {
		t.Errorf("Format(fallback) = %q; want %q", got, want)
	}
}

func TestFormat_FallbackNoShowName(t *testing.T) {
	cfg := ConfigSection{FallbackIcon: "?", ShowName: false}
	entry := IconEntry{}

	got := Format(entry, cfg, "zsh")
	want := "?"
	if got != want {
		t.Errorf("Format(fallback) = %q; want %q", got, want)
	}
}

func TestMerge_UserOverridesEmbedded(t *testing.T) {
	embedded := &Config{
		Config: ConfigSection{FallbackIcon: "?", ShowName: true},
		Icons: map[string]IconEntry{
			"nvim": {DisplayName: "nvim", Icon: "N"},
		},
	}
	user := &Config{
		Config: ConfigSection{FallbackIcon: "X"},
		Icons: map[string]IconEntry{
			"nvim": {DisplayName: "nvim", Icon: "X"},
		},
	}

	got := merge(user, embedded)

	if got.Config.FallbackIcon != "X" {
		t.Errorf("merge: FallbackIcon = %q; want X", got.Config.FallbackIcon)
	}
	if got.Icons["nvim"].Icon != "X" {
		t.Errorf("merge: Icons[nvim].Icon = %q; want X", got.Icons["nvim"].Icon)
	}
}

func TestMerge_EmbeddedEntriesPreserved(t *testing.T) {
	embedded := &Config{
		Config: ConfigSection{FallbackIcon: "?"},
		Icons: map[string]IconEntry{
			"nvim": {DisplayName: "nvim", Icon: "N"},
			"git":  {DisplayName: "git", Icon: "G"},
		},
	}
	user := &Config{
		Config: ConfigSection{FallbackIcon: "X"},
		Icons: map[string]IconEntry{
			"nvim": {DisplayName: "nvim", Icon: "X"},
		},
	}

	got := merge(user, embedded)

	if _, ok := got.Icons["git"]; !ok {
		t.Errorf("merge: embedded entry git not preserved")
	}
	if got.Icons["nvim"].Icon != "X" {
		t.Errorf("merge: user entry nvim not overriding")
	}
}

func TestParseAndValidate_MissingIconAndPattern(t *testing.T) {
	data := []byte(`
icons:
  bad:
    icon: ""
`)
	_, err := parseAndValidate(data)
	if err != ErrNoIconOrPattern {
		t.Errorf("parseAndValidate() error = %v; want ErrNoIconOrPattern", err)
	}
}

func TestParseAndValidate_InvalidRegex(t *testing.T) {
	data := []byte(`
icons:
  bad:
    icon: "X"
    pattern: "[invalid"
`)
	_, err := parseAndValidate(data)
	if err == nil {
		t.Errorf("parseAndValidate() error = nil; want invalid regex error")
	}
}

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Config.FallbackIcon != "?" {
		t.Errorf("FallbackIcon = %q; want \"?\"", cfg.Config.FallbackIcon)
	}

	entry, ok := cfg.Icons["nvim"]
	if !ok {
		t.Fatalf("Icons[nvim] not found")
	}
	if entry.Icon == "" {
		t.Errorf("Icons[nvim].Icon = %q; want non-empty", entry.Icon)
	}
	if entry.Pattern != nil {
		t.Errorf("Icons[nvim].Pattern = %v; want nil", entry.Pattern)
	}
}

func TestParseAndValidate_ValidEntryWithEmptyIcon(t *testing.T) {
	data := []byte(`
icons:
  db:
    icon: ""
    pattern: "nvim"
`)
	cfg, err := parseAndValidate(data)
	if err != nil {
		t.Errorf("parseAndValidate() error = %v; want nil", err)
	}
	if cfg.Icons["db"].Icon != "" || cfg.Icons["db"].Pattern == nil {
		t.Errorf("parseAndValidate() = %+v; want icon='', pattern=compiled", cfg.Icons["db"])
	}
}
