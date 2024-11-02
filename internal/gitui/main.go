package gitui

import (
	"embed"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed gitui
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `gitui`,
	Usage: `gitui <command>`,
	Short: `gitui is a utility to manage gitui configuration`,
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
	Short: `Setup gitui`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "gitui", oscfg.ConfigDir())
	},
}
