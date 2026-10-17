package tmux

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installTmux() error {
	if ok, _ := have.Executable(`tmux`); ok {
		fmt.Println(`tmux is already installed`)
		return nil
	}
	switch distro.ID() {
	case `arch`:
		return install.WithRoot(`pacman`, `-S`, `--noconfirm`, `tmux`)
	case `ubuntu`, `debian`:
		return install.WithRoot(`apt-get`, `install`, `-y`, `tmux`)
	case `fedora`:
		return run.Exec(`dnf`, `install`, `tmux`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `tmux`)
	default:
		fmt.Fprintln(os.Stderr, `Unsupported operating system. Please install tmux manually.`)
	}
	return nil
}
