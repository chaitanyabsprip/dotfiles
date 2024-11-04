package distro

import (
	"fmt"

	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai"
)

var Cmd = &bonzai.Cmd{
	Name:  `distro`,
	Usage: `distro <command>`,
	Short: `distro <command>`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
	Call: func(x *bonzai.Cmd, args ...string) error {
		fmt.Println(Name())
		return nil
	},
}

