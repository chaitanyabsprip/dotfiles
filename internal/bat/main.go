package bat

import (
	"embed"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
	"github.com/Chaitanyabsprip/dot/x/install"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed bat
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bat`,
	Short: `bat is a utility to manage bat configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, install.BatCmd.WithName(`install`)},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `Setup bat`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, `bat`, oscfg.ConfigDir(), nil)
	},
}
