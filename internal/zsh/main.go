package zsh

import (
	"embed"
	"os"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed all:zsh
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `zsh`,
	Usage: `zsh <command>`,
	Short: `zsh is a utility to manage zsh configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
	Init: func(x *bonzai.Cmd, args ...string) error {
		for _, cmd := range x.Cmds {
			cmd.Caller = x
		}
		return nil
	},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Usage: `setup <opts>`,
	Opts:  `slim|quik|full`,
	Short: `Setup zsh`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(
			embedFs,
			`zsh`,
			oscfg.ConfigDir(),
			map[string]string{
				`zsh/.zshenv`: os.Getenv("HOME"),
			},
		)
	},
}
