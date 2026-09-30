// Package neovim installs neovim from the system package manager and
// clones its config from a separate repo — unlike every other tool,
// neovim's config is not embedded. It lives at
// https://github.com/Chaitanyabsprip/nvim as its own repo (a git
// submodule in the real dotfiles repo), too large and independently
// versioned to vendor into this binary the way every other tool's
// config is.
package neovim

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

const repoURL = `https://github.com/Chaitanyabsprip/nvim.git`

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Short: `install neovim`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		if ok, _ := have.Executable(`nvim`); ok {
			fmt.Println(`nvim is already installed`)
			return nil
		}
		// Debian's neovim is too old for this config; build from source there.
		return install.Pkg(`neovim`, map[string]string{`apt`: ``})
	},
}

var SetupCmd = &bonzai.Cmd{
	Name: `setup`,

	Short: `clone the neovim config into ~/.config/nvim`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		dest := filepath.Join(oscfg.ConfigDir(), `nvim`)
		if !futil.Exists(dest) {
			return run.Exec(`git`, `clone`, `--branch`, `main`, repoURL, dest)
		}
		if !futil.Exists(filepath.Join(dest, `.git`)) {
			return fmt.Errorf(`%s exists and isn't a git repo — resolve manually`, dest)
		}
		// --ff-only: never clobber local commits/changes, just fail loudly
		// if history diverged, unlike every other tool's drift-aware copy.
		return run.Exec(`git`, `-C`, dest, `pull`, `--ff-only`)
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit neovim configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(oscfg.ConfigDir(), `nvim`, `init.lua`)
		return edit.Files(filePath)
	},
}
