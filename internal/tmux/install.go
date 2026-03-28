package tmux

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

func installTmux() error {
	if ok, _ := have.Executable(`tmux`); ok {
		fmt.Println(`tmux is already installed`)
		return nil
	}
	switch distro.Name() {
	case `Arch Linux`:
		return withRoot(`pacman`, `-S`, `tmux`)
	case `Ubuntu`, `Debian GNU/Linux`:
		return withRoot(`apt-get`, `install`, `-y`, `tmux`)
	case `Fedora Linux`:
		return run.Exec(`dnf`, `install`, `tmux`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `tmux`)
	default:
		fmt.Fprintln(os.Stderr, `Unsupported operating system. Please install tmux manually.`)
	}
	return nil
}

func withRoot(args ...string) error {
	if os.Geteuid() != 0 {
		if _, err := exec.LookPath(`sudo`); err != nil {
			return fmt.Errorf(`user not root and sudo not found`)
		}
		args = append([]string{`sudo`}, args...)
	}
	return run.Exec(args...)
}
