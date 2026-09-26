// Package neovim installs neovim from the system package manager.
package neovim

import (
	"fmt"

	"github.com/rwxrob/bonzai"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

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
