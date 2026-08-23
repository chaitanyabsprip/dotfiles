// Package dot provides dotfile management functionality through a CLI application.
// It allows setting up, installing, and managing various configuration files for
// different tools and applications across the system.
package dot

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/claude"
	idot "github.com/Chaitanyabsprip/dotfiles/internal/dot"
	"github.com/Chaitanyabsprip/dotfiles/x"
)

var Cmd = &bonzai.Cmd{
	Name:  `dot`,
	Alias: `d`,
	Short: `manage dotfiles`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		idot.SetupCmd,
		idot.InstallCmd,
		idot.EditCmd,
		idot.DepsCmd,
		idot.InitCmd,
		x.Cmd,
		claude.Cmd,
		help.Cmd,
	},
}
