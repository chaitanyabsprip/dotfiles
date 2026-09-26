// Package testenv sandboxes tests from the user's home directory and
// tmux server. Call Isolate from a package's TestMain.
package testenv

import (
	"os"
	"path/filepath"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

// Isolate points HOME, the XDG base dirs and the dotfiles path variables
// at a fresh temp dir, in both the process environment and pkg/env (which
// read them once at startup). It unsets TMUX and points TMUX_TMPDIR at
// the temp dir, so a tmux call finds no server instead of the user's.
// cleanup restores everything and removes the temp dir.
func Isolate() (home string, cleanup func(), err error) {
	home, err = os.MkdirTemp(``, `dotfiles-test-home-`)
	if err != nil {
		return ``, nil, err
	}
	if home, err = filepath.EvalSymlinks(home); err != nil {
		return ``, nil, err
	}
	vars := map[string]string{
		`HOME`:            home,
		`XDG_CONFIG_HOME`: filepath.Join(home, `.config`),
		`XDG_CACHE_HOME`:  filepath.Join(home, `.cache`),
		`XDG_STATE_HOME`:  filepath.Join(home, `.local`, `state`),
		`XDG_DATA_HOME`:   filepath.Join(home, `.local`, `share`),
		`SCRIPTS`:         filepath.Join(home, `.local`, `bin`),
		`PROJECTS`:        filepath.Join(home, `projects`),
		`PROGRAMS`:        filepath.Join(home, `programs`),
		`PICTURES`:        filepath.Join(home, `pictures`),
		`DOWNLOADS`:       filepath.Join(home, `downloads`),
		`NOTESPATH`:       filepath.Join(home, `notes`),
		`DOTFILES`:        filepath.Join(home, `dotfiles`),
		`TMUX_TMPDIR`:     filepath.Join(home, `tmux`),
		`TMUX`:            ``,
	}
	saved := map[string]*string{}
	for k, v := range vars {
		if old, ok := os.LookupEnv(k); ok {
			saved[k] = &old
		} else {
			saved[k] = nil
		}
		if v == `` {
			os.Unsetenv(k)
		} else {
			os.Setenv(k, v)
		}
	}
	savedPkg := snapshot()
	loadPkg()
	return home, func() {
		for k, old := range saved {
			if old == nil {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, *old)
			}
		}
		savedPkg.restore()
		os.RemoveAll(home)
	}, nil
}

type pkgVars struct {
	home, cache, config, state, scripts, projects, programs string
	pictures, downloads, notes, dotfiles, tmux              string
}

func snapshot() pkgVars {
	return pkgVars{
		env.Home, env.XdfCacheHome, env.XdgConfigHome, env.XdgStateHome,
		env.Scripts, env.Projects, env.Programs, env.Pictures,
		env.Downloads, env.Notespath, env.Dotfiles, env.Tmux,
	}
}

func (p pkgVars) restore() {
	env.Home, env.XdfCacheHome, env.XdgConfigHome, env.XdgStateHome = p.home, p.cache, p.config, p.state
	env.Scripts, env.Projects, env.Programs, env.Pictures = p.scripts, p.projects, p.programs, p.pictures
	env.Downloads, env.Notespath, env.Dotfiles, env.Tmux = p.downloads, p.notes, p.dotfiles, p.tmux
}

// loadPkg re-reads pkg/env's copies from the (now isolated) environment.
func loadPkg() {
	pkgVars{
		os.Getenv(`HOME`), os.Getenv(`XDG_CACHE_HOME`), os.Getenv(`XDG_CONFIG_HOME`),
		os.Getenv(`XDG_STATE_HOME`), os.Getenv(`SCRIPTS`), os.Getenv(`PROJECTS`),
		os.Getenv(`PROGRAMS`), os.Getenv(`PICTURES`), os.Getenv(`DOWNLOADS`),
		os.Getenv(`NOTESPATH`), os.Getenv(`DOTFILES`), os.Getenv(`TMUX`),
	}.restore()
}
