// Package ripgrep installs ripgrep (rg).
package ripgrep

import (
	"fmt"

	"github.com/rwxrob/bonzai"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Short: `install ripgrep`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		if ok, _ := have.Executable(`rg`); ok {
			fmt.Println(`rg is already installed`)
			return nil
		}
		return install.Pkg(`ripgrep`, nil)
	},
}
