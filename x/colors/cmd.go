// Package colors provides commands for printing and manipulating terminal colors
// in different formats. It includes utilities to display color tables, strip color
// codes, and visualize terminal color capabilities.
package colors

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/sunrise"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `color`,
	Short: `print colors in terminal in different formats`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		Color255Cmd,
		TableCmd,
		StripCmd,
		sunrise.Cmd,
		TermCmd,
	},
}
