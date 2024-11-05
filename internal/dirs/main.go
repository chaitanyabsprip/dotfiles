package dirs

import (
	"embed"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed dirs
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `dirs`,
	Usage: `dirs <command>`,
	Short: `dirs is a utility to manage dirs configuration`,
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
	Short: `Setup dirs`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "dirs", oscfg.ConfigDir(),
			map[string]string{
				`dirs`: filepath.Join(
					oscfg.ConfigDir(),
					`user-dirs.dirs`,
				),
			},
		)
	},
}
