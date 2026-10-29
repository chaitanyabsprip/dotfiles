package zsh

import (
	"fmt"
	"os"

	"github.com/chaitanyabsprip/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installZsh() error {
	if ok, _ := have.Executable(`zsh`); ok {
		fmt.Println(`zsh is already installed`)
		return nil
	}
	switch distro.ID() {
	case `arch`:
		return install.WithRoot(`pacman`, `-S`, `--noconfirm`, `zsh`)
	case `ubuntu`, `debian`:
		return install.WithRoot(`apt-get`, `install`, `-y`, `zsh`)
	case `fedora`:
		return run.Exec(`dnf`, `install`, `zsh`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `zsh`)
	default:
		fmt.Fprintln(os.Stderr, `Unsupported operating system. Please install zsh manually.`)
	}
	return nil
}
