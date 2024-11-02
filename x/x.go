package x

import (
	"github.com/rwxrob/bonzai/comp"
	bonzai "github.com/rwxrob/bonzai/pkg"

	"github.com/Chaitanyabsprip/dot/x/colors"
	"github.com/Chaitanyabsprip/dot/x/have"
	"github.com/Chaitanyabsprip/dot/x/workdirs"
)

var Cmd = &bonzai.Cmd{
	Name:  `x`,
	Usage: `x <command>`,
	Short: `x is a utility to manage scripts`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		have.Cmd,
		workdirs.Cmd,
		colors.Cmd,
		creashCmd,
	},
}
