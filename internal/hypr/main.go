package hypr

import (
	"embed"
	"runtime"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed hypr
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `hypr`,
	Short: `manage hypr configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup hypr`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if runtime.GOOS == "linux" {
			return nil
		}
		return e.SetupAll(embedFs, `hypr`, oscfg.ConfigDir(), nil)
	},
}
