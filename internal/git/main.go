package git

import (
	"embed"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/comp"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed git
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `git`,
	Usage: `git <command>`,
	Short: `git is a utility to manage git configuration`,
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
	Short: `Setup git`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "git", oscfg.ConfigDir())
	},
}
