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
	switch distro.Name() {
	case `Arch Linux`:
		return install.WithRoot(`pacman`, `-S`, `tmux`)
	case `Ubuntu`, `Debian GNU/Linux`:
		return install.WithRoot(`apt-get`, `install`, `-y`, `tmux`)
	case `Fedora Linux`:
		return run.Exec(`dnf`, `install`, `tmux`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `tmux`)
	default:
		fmt.Fprintln(os.Stderr, `Unsupported operating system. Please install tmux manually.`)
	}
	return nil
}
