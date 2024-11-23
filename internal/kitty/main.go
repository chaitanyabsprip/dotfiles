package kitty

import (
	"embed"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed kitty
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `kitty`,
	Short: `kitty is a utility to manage kitty configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup kitty`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "kitty", oscfg.ConfigDir(), nil)
	},
}
