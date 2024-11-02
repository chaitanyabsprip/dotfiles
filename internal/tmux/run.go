package tmux

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/comp"
)

var runCmd = &bonzai.Cmd{
	Name:  `run`,
	Alias: `x`,
	Usage: `tmux x`,
	Short: `tmux x`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
}

