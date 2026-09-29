package say

import (
	"fmt"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `say`,
	Short: `print a short colored status message`,
	Long: `
Prints ` + "`msg`" + ` in one of six styles: error, warning, success,
inprogress, bold or italic. See 'say help' for each.`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{errorCmd, warningCmd, successCmd, inProgressCmd, boldCmd, italicCmd, help.Cmd},
	Do: func(x *bonzai.Cmd, _ ...string) error {
		fmt.Printf("%s - %s\n\n", x.Name, x.Short)
		fmt.Println(`COMMANDS:`)
		for _, c := range x.Cmds {
			fmt.Printf("  %-10s - %s\n", c.Name, c.Short)
		}
		return nil
	},
}

var errorCmd = leaf(`error`, `print msg in red with an error mark`, `red, with an error mark`, Error)
var warningCmd = leaf(`warning`, `print msg in yellow with a warning mark`, `yellow, with a warning mark`, Warning)
var successCmd = leaf(`success`, `print msg in green with a checkmark`, `green, with a checkmark`, Success)
var inProgressCmd = leaf(`inprogress`, `print msg in blue prefixed with '...'`, `blue, prefixed with '...'`, InProgress)
var boldCmd = leaf(`bold`, `print msg in bold`, `bold`, Bold)
var italicCmd = leaf(`italic`, `print msg in italic`, `italic`, Italic)

// leaf builds a command that joins its args into msg and prints
// style(msg).
func leaf(name, short, styleDesc string, style func(msg string) string) *bonzai.Cmd {
	return &bonzai.Cmd{
		Name:  name,
		Short: short,
		Long: `
Prints ` + "`msg`" + ` (its args joined with spaces), ` + styleDesc + `.`,
		Do: func(_ *bonzai.Cmd, args ...string) error {
			fmt.Println(style(strings.Join(args, ` `)))
			return nil
		},
	}
}
