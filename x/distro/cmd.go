// Package distro provides utilities for identifying and querying information about
// the operating system distribution. It can detect various Linux distributions and
// other Unix-like systems, and provides methods to retrieve their names and versions.
package distro

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `distro`,
	Short: `distro <command>`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
	Do: func(x *bonzai.Cmd, args ...string) error {
		fmt.Println(Name())
		return nil
	},
}
