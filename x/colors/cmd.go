package colors

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/sunrise"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `color`,
	Short: `Print colors in terminal in different formats`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		Color255Cmd,
		TableCmd,
		StripCmd,
		sunrise.Cmd,
		TermCmd,
	},
}
