package colors

import (
	"github.com/rwxrob/bonzai/comp"
	bonzai "github.com/rwxrob/bonzai/pkg"
)

var Cmd = &bonzai.Cmd{
	Name:  `color`,
	Usage: `color <command>`,
	Short: `Print colors in terminal in different formats`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		Color255Cmd,
		TableCmd,
		StripCmd,
		SunriseCmd,
		TermCmd,
	},
}
