package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/futil"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

var OhMyPoshCmd = &bonzai.Cmd{
	Name: `ohmyposh`,
	Do: func(x *bonzai.Cmd, args ...string) error { return OhMyPosh() },
}

func OhMyPosh() error {
	if ok, _ := have.Executable(`ohmyposh`); ok {
		fmt.Println(`ohmyposh is already installed`)
		return nil
	}
	binDir := filepath.Join(os.Getenv(`HOME`), `.local`, `bin`)
	if err := futil.CreateDir(binDir); err != nil {
		return err
	}
	if err := unzipPkgInstall(); err != nil {
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
	if err := bash.Run(); err != nil {
		return err
	}
	return nil
}
