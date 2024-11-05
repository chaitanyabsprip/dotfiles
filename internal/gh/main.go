package gh

import (
	"embed"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

//go:embed gh
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `gh`,
	Usage: `gh <command>`,
	Short: `gh is a utility to manage github-cli configuration`,
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
	Short: `Setup github-cli`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "gh", oscfg.ConfigDir(), nil)
	},
}
