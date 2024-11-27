package install

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

var ZshCmd = &bonzai.Cmd{
	Name:  `zsh`,
	Alias: `i`,
	Do:    func(_ *bonzai.Cmd, _ ...string) error { return Zsh() },
}

func Zsh() error {
	if ok, _ := have.Executable(`zsh`); ok {
		fmt.Println(`zsh is already installed`)
		return nil
	}
	switch distro.Name() {
	case `Arch Linux`:
		return WithRoot(`pacman`, `-S`, `zsh`)
	case `Ubuntu`, `Debian GNU/Linux`:
		return WithRoot(`apt-get`, `install`, `-y`, `zsh`)
	case `Fedora Linux`:
		return run.Exec(`dnf`, `install`, `zsh`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `zsh`)
	default:
		fmt.Fprintln(
			os.Stderr,
			`Unsupported operating system. Please install zsh manually.`,
		)
	}
	return nil
}

var ZapCmd = &bonzai.Cmd{
	Name: `zap`,
	Do:   func(*bonzai.Cmd, ...string) error { return Zap() },
}

func Zap() error {
	if path, err := exec.LookPath(`zap`); err == nil {
		fmt.Println(`zap is already installed at`, path)
		return err
	}
	err := DownloadFile(
		`https://raw.githubusercontent.com/zap-zsh/zap/master/install.zsh`,
		`install.zsh`,
	)
	if err != nil {
		return err
	}
	err = run.Exec(
		`zsh`,
		`install.zsh`,
		`--branch`,
		`release-v1`,
		`--keep`,
	)
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
