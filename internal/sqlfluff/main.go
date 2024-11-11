package sqlfluff

import (
	"embed"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed all:sqlfluff
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `sqlfluff`,
	Usage: `sqlfluff <command>`,
	Short: `sqlfluff is a utility to manage sqlfluff configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Usage: `setup <opts>`,
	Opts:  `slim|quik|full`,
	Short: `Setup sqlfluff`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "sqlfluff", oscfg.ConfigDir(), nil)
	},
}
