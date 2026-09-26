// Package yq installs mikefarah/yq.
package yq

import (
	"fmt"

	"github.com/rwxrob/bonzai"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Short: `install yq`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		if ok, _ := have.Executable(`yq`); ok {
			fmt.Println(`yq is already installed`)
			return nil
		}
		if install.IsApt() { // Debian's yq package is a different tool
			return install.GhRelease(`mikefarah/yq`, `yq`, yqAsset)
		}
		return install.Pkg(`yq`, map[string]string{`pacman`: `go-yq`})
	},
}

func yqAsset(_, goos, goarch string) string {
	if (goos != `linux` && goos != `darwin`) || (goarch != `amd64` && goarch != `arm64`) {
		return ``
	}
	return `yq_` + goos + `_` + goarch
}
