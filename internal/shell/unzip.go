package shell

import (
	"fmt"

	"github.com/chaitanyabsprip/bonzai"
	"github.com/chaitanyabsprip/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

var InstallUnzipCmd = &bonzai.Cmd{
	Name: `unzip`,
	Do:   func(x *bonzai.Cmd, args ...string) error { return InstallUnzip() },
}

func InstallUnzip() error {
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
	switch distro.ID() {
	case `arch`:
		return install.WithRoot(`pacman`, `-S`, `--noconfirm`, `unzip`)
	case `ubuntu`, `debian`:
		return install.WithRoot(`apt-get`, `install`, `unzip`, `-y`)
	case `fedora`:
		return run.Exec(`dnf`, `install`, `unzip`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `unzip`)
	default:
		return fmt.Errorf(
			`unsupported or unconfigured operating system`,
		)
	}
}
