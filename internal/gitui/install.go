package gitui

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installGitui() error {
	if ok, _ := have.Executable(`gitui`); ok {
		fmt.Println(`gitui is already installed`)
		return nil
	}
	if install.IsApt() { // no gitui package on Debian or Ubuntu
		return install.GhRelease(`gitui-org/gitui`, `gitui`, gituiAsset)
	}
	return install.Pkg(`gitui`, nil)
}

func gituiAsset(_, goos, goarch string) string {
	switch goos + `/` + goarch {
	case `linux/amd64`:
		return `gitui-linux-x86_64.tar.gz`
	case `linux/arm64`:
		return `gitui-linux-aarch64.tar.gz`
	case `darwin/arm64`:
		return `gitui-mac.tar.gz`
	case `darwin/amd64`:
		return `gitui-mac-x86.tar.gz`
	}
	return ``
}
