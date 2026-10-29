// Package dot provides dotfile management functionality through a CLI application.
// It allows setting up, installing, and managing various configuration files for
// different tools and applications across the system.
package dot

import (
	"github.com/chaitanyabsprip/bonzai"
	"github.com/chaitanyabsprip/bonzai/cmds/help"
	"github.com/chaitanyabsprip/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/cmdhelp"
	"github.com/Chaitanyabsprip/dotfiles/internal/core/version"
	idot "github.com/Chaitanyabsprip/dotfiles/internal/dot"
	"github.com/Chaitanyabsprip/dotfiles/x"
)

var Cmd = &bonzai.Cmd{
	Name:  `dot`,
	Alias: `d`,
	Short: `manage dotfiles`,
	Comp:  comp.Cmds,
	Funcs: cmdhelp.Funcs,
	Cmds: []*bonzai.Cmd{
		idot.SetupCmd,
		idot.InstallCmd,
		idot.EditCmd,
		idot.DepsCmd,
		idot.StatusCmd,
		idot.DiffCmd,
		idot.InitCmd,
		x.Cmd,
		version.Cmd,
		help.Cmd,
	},
}
