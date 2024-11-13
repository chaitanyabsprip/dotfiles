package zsh

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/run"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
	"github.com/Chaitanyabsprip/dot/x/install"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed all:zsh
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `zsh`,
	Short: `zsh is a utility to manage zsh configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `Setup zsh`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		zshenvPath := filepath.Join(os.Getenv(`HOME`), `.zshenv`)
		err := e.SetupAll(
			embedFs,
			`zsh`,
			oscfg.ConfigDir(),
			map[string]string{`zsh/.zshenv`: zshenvPath},
		)
		if err != nil {
			return err
		}
		return installZap()
	},
}

func installZap() error {
	if path, err := exec.LookPath(`zap`); err == nil {
		fmt.Println(`zap is already installed at`, path)
		return err
	}
	err := install.DownloadFile(
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
