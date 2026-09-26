// Package eza installs eza, a modern ls.
package eza

import (
	"fmt"

	"github.com/rwxrob/bonzai"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Short: `install eza`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		if ok, _ := have.Executable(`eza`); ok {
			fmt.Println(`eza is already installed`)
			return nil
		}
		if install.IsApt() { // older Debian and Ubuntu have no eza package
			return install.GhRelease(`eza-community/eza`, `eza`, ezaAsset)
		}
		return install.Pkg(`eza`, nil)
	},
}

func ezaAsset(_, goos, goarch string) string {
	if goos != `linux` {
		return ``
	}
	switch goarch {
	case `amd64`:
		return `eza_x86_64-unknown-linux-gnu.tar.gz`
	case `arm64`:
		return `eza_aarch64-unknown-linux-gnu.tar.gz`
	}
	return ``
}
