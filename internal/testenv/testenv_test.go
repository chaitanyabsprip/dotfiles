package testenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

func TestIsolate(t *testing.T) {
	realHome := os.Getenv(`HOME`)
	t.Setenv(`TMUX`, `/tmp/tmux-1/default,1,0`) // pretend we run inside tmux

	home, cleanup, err := Isolate()
	if err != nil {
		t.Fatal(err)
	}
	if home == realHome || !filepath.IsAbs(home) {
		t.Fatalf("home = %q, want a fresh absolute temp dir", home)
	}
	for _, k := range []string{`HOME`, `XDG_CONFIG_HOME`, `XDG_CACHE_HOME`, `XDG_STATE_HOME`,
		`XDG_DATA_HOME`, `SCRIPTS`, `PROJECTS`, `TMUX_TMPDIR`} {
		if v := os.Getenv(k); !strings.HasPrefix(v, home) {
			t.Errorf("%s = %q, want it under %s", k, v, home)
		}
	}
	if env.Home != home || !strings.HasPrefix(env.XdfCacheHome, home) || !strings.HasPrefix(env.Scripts, home) {
		t.Errorf("pkg/env not isolated: Home=%q Cache=%q Scripts=%q", env.Home, env.XdfCacheHome, env.Scripts)
	}
	if os.Getenv(`TMUX`) != `` || env.Tmux != `` {
		t.Error(`TMUX should be unset so tests cannot reach the user's tmux server`)
	}

	cleanup()
	if _, err := os.Stat(home); !os.IsNotExist(err) {
		t.Errorf("temp home not removed: %v", err)
	}
	if os.Getenv(`HOME`) != realHome || env.Home != realHome {
		t.Errorf("HOME not restored: env=%q pkg=%q", os.Getenv(`HOME`), env.Home)
	}
}
