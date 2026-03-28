package bash

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

func installOhMyPosh() error {
	if ok, _ := have.Executable(`ohmyposh`); ok {
		fmt.Println(`ohmyposh is already installed`)
		return nil
	}
	binDir := oscfg.BinDir()
	if !futil.Exists(binDir) {
		if err := futil.CreateDir(binDir); err != nil {
			return err
		}
	}
	if err := installUnzip(); err != nil {
		return err
	}
	cmd := exec.Command(`curl`, `-s`, `https://ohmyposh.dev/install.sh`)
	bash := exec.Command(`bash`, `-s`)
	bash.Stdin, _ = cmd.StdoutPipe()
	bash.Stdout = os.Stdout
	bash.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return bash.Run()
}

func installUnzip() error {
	switch distro.Name() {
	case `Arch Linux`:
		return withRoot(`pacman`, `-S`, `unzip`)
	case `Ubuntu`, `Debian GNU/Linux`:
		return withRoot(`apt-get`, `install`, `unzip`, `-y`)
	case `Fedora`:
		return run.Exec(`dnf`, `install`, `unzip`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `unzip`)
	default:
		return fmt.Errorf(`unsupported or unconfigured operating system`)
	}
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
