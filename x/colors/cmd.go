// Package colors provides commands for printing and manipulating terminal colors
// in different formats. It includes utilities to display color tables, strip color
// codes, and visualize terminal color capabilities.
package colors

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/cmds/sunrise"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `color`,
	Short: `print colors in terminal in different formats`,
	Long: `
Prints color tables and demos in different formats. See 'color help'
for the individual commands.`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{
		Color255Cmd,
		TableCmd,
		StripCmd,
		sunrise.Cmd,
		TermCmd,
		help.Cmd,
	},
	Do: func(x *bonzai.Cmd, _ ...string) error {
		fmt.Printf("%s - %s\n\n", x.Name, x.Short)
		fmt.Println(`COMMANDS:`)
		for _, c := range x.Cmds {
			fmt.Printf("  %-10s - %s\n", c.Name, c.Short)
		}
		return nil
	},
}
