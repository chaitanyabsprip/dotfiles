package tmux

import (
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai"
)

var runCmd = &bonzai.Cmd{
	Name:  `run`,
	Alias: `x`,
	Usage: `tmux x`,
	Short: `tmux x`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
}
