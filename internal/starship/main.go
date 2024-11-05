package starship

import (
	"embed"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed starship.toml
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `starship`,
	Usage: `starship <command>`,
	Short: `starship is a utility to manage starship configuration`,
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
	Short: `Setup starship`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "starship", oscfg.ConfigDir(), nil)
	},
}
