package brew

import (
	"embed"
	"runtime"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed Brewfile
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `brew`,
	Short: `brew is a utility to manage brew configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `Setup brew`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if runtime.GOOS == "darwin" {
			return nil
		}
		return e.SetupAll(embedFs, "brew", oscfg.ConfigDir(), nil)
	},
}
