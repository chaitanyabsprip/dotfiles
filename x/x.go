// Package x is a utility to manage scripts and provide various command-line tools
// for system configuration, development, and productivity. It serves as the main
// entry point for several subcommands that perform specific functions.
package x

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/x/colors"
	"github.com/Chaitanyabsprip/dotfiles/x/depends"
	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/gpt"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
	"github.com/Chaitanyabsprip/dotfiles/x/workdirs"
)

var Cmd = &bonzai.Cmd{
	Name:  `x`,
	Short: `x is a utility to manage scripts`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		colors.Cmd,
		catcCmd,
		creashCmd,
		depends.Cmd,
		distro.Cmd,
		gpt.Cmd,
		have.Cmd,
		help.Cmd,
		install.Cmd,
		workdirs.Cmd,
	},
}
