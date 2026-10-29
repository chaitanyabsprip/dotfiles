// Package distro provides utilities for identifying and querying information about
// the operating system distribution. It can detect various Linux distributions and
// other Unix-like systems, and provides methods to retrieve their names and versions.
package distro

import (
	"fmt"

	"github.com/chaitanyabsprip/bonzai"
	"github.com/chaitanyabsprip/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `distro`,
	Short: `print the OS distro id and version`,
	Long: `
Prints "id:version" for the current OS (e.g. "ubuntu:24.04"): ID and
VERSION_ID from /etc/os-release on Linux, from lsb_release if that file
is missing, or from uname otherwise (e.g. "Darwin:25.5.0" on macOS).`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{},
	Do: func(x *bonzai.Cmd, args ...string) error {
		fmt.Printf("%s:%s\n", ID(), Version())
		return nil
	},
}
