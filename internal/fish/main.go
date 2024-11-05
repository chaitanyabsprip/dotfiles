package fish

import (
	"embed"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed fish
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `fish`,
	Usage: `fish <command>`,
	Short: `fish is a utility to manage fish configuration`,
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
	Short: `Setup fish`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "fish", oscfg.ConfigDir(), nil)
	},
}
