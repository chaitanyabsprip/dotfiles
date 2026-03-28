package zsh

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/web"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/pkg/with"
	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

func installZsh() error {
	if ok, _ := have.Executable(`zsh`); ok {
		fmt.Println(`zsh is already installed`)
		return nil
	}
	switch distro.Name() {
	case `Arch Linux`:
		return withRoot(`pacman`, `-S`, `zsh`)
	case `Ubuntu`, `Debian GNU/Linux`:
		return withRoot(`apt-get`, `install`, `-y`, `zsh`)
	case `Fedora Linux`:
		return run.Exec(`dnf`, `install`, `zsh`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `zsh`)
	default:
		fmt.Fprintln(os.Stderr, `Unsupported operating system. Please install zsh manually.`)
	}
	return nil
}

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

func installZap() (err error) {
	if path, err := exec.LookPath(`zap`); err == nil {
		fmt.Println(`zap is already installed at`, path)
		return err
	}
	zshDir := filepath.Join(oscfg.ConfigDir(), `zsh`)
	pop, err := with.Env(`ZDOTDIR`, zshDir)
	if err != nil {
		return err
	}
	defer func() { err = pop() }()
	err = downloadFile(
		`https://raw.githubusercontent.com/zap-zsh/zap/master/install.zsh`,
		`install.zsh`,
	)
	if err != nil {
		return err
	}
	err = run.Exec(`zsh`, `install.zsh`, `--branch`, `release-v1`, `--keep`)
	if err != nil {
		return err
	}
	err = os.Remove(`install.zsh`)
	if err != nil {
		return err
	}
	fmt.Println(`zap installed`)
	return nil
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

func downloadFile(url, dest string) (err error) {
	file, err := os.Create(dest)
	if err != nil {
		return
	}
	defer func() { err = errors.Join(err, file.Close()) }()
	req := web.Req{U: url, D: file}
	err = req.Submit()
	return
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
