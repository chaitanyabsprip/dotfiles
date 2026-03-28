// Package dot provides dependency management infrastructure for dotfile tools.
package dot

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

// DepsCmd is infrastructure for tools to declare their dependencies.
// Individual tool packages can register subcommands here in the future.
var DepsCmd = &bonzai.Cmd{
	Name:  `deps`,
	Short: `check dependencies for a given tool`,
	Comp:  comp.Cmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: deps <tool> — specify a tool to check its dependencies")
		}
		fmt.Printf("%s: no dependencies defined\n", args[0])
		return nil
	},
}
