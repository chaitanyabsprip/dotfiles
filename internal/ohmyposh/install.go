package ohmyposh

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/futil"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/internal/shell"
	"github.com/Chaitanyabsprip/dotfiles/x/depends"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

var InstallCmd = &bonzai.Cmd{
	Name: `ohmyposh`,
	Do:   func(x *bonzai.Cmd, args ...string) error { return installOhMyPosh() },
}

// Deps lists ohmyposh's real runtime dependencies, per ADR-0011 — see
// installOhMyPosh below, which shells out to unzip.
var Deps = []depends.Dep{
	{Name: `unzip`, Note: `install`},
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
	if err := shell.InstallUnzip(); err != nil {
		return err
	}
	depends.On(func(err error) {
		fmt.Println(err.Error())
	}, `unzip`)
	cmd := exec.Command(`curl`, `-s`, `https://ohmyposh.dev/install.sh`)
	bash := exec.Command(`bash`, `-s`)
	bash.Stdin, _ = cmd.StdoutPipe()
	bash.Stdout = os.Stdout
	bash.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := bash.Run(); err != nil {
		return err
	}
	return nil
}
