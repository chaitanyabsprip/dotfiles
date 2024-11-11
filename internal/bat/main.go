package bat

import (
	"embed"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed bat
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bat`,
	Usage: `bat <command>`,
	Short: `bat is a utility to manage bat configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Usage: `setup <opts>`,
	Opts:  `slim|quik|full`,
	Short: `Setup bat`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "bat", oscfg.ConfigDir(), nil)
	},
}
