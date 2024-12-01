package install

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

var UnzipCmd = &bonzai.Cmd{
	Name: `unzip`,
	Do:   func(x *bonzai.Cmd, args ...string) error { return Unzip() },
}

func Unzip() error {
	if ok, _ := have.Executable(`unzip`); ok {
		fmt.Println(`unzip is already installed`)
		return nil
	}
	if err := unzipPkgInstall(); err != nil {
		return err
	}
	return nil
}

func unzipPkgInstall() error {
	fmt.Println(distro.Name())
	switch distro.Name() {
	case `Arch Linux`:
		return WithRoot(`pacman`, `-S`, `unzip`)
	case `Ubuntu`, `Debian GNU/Linux`:
		return WithRoot(`apt-get`, `install`, `unzip`, `-y`)
	case `Fedora`:
		return run.Exec(`dnf`, `install`, `unzip`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `unzip`)
	default:
		return fmt.Errorf(
			`unsupported or unconfigured operating system`,
		)
	}
}
