package alacritty

import (
	"embed"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed alacritty
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `alacritty`,
	Usage: `alacritty <command>`,
	Short: `alacritty is a utility to manage alacritty configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Usage: `setup <opts>`,
	Opts:  `slim|quik|full`,
	Short: `Setup alacritty`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "alacritty", oscfg.ConfigDir(), nil)
	},
}
