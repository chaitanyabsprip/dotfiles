// Package jq installs jq.
package jq

import (
	"fmt"

	"github.com/rwxrob/bonzai"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Short: `install jq`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		if ok, _ := have.Executable(`jq`); ok {
			fmt.Println(`jq is already installed`)
			return nil
		}
		return install.Pkg(`jq`, nil)
	},
}
