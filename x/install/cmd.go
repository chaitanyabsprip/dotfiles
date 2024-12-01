package install

import (
	"github.com/rwxrob/bonzai/comp"

	"github.com/rwxrob/bonzai"
)

var Cmd = &bonzai.Cmd{
	Name: `install`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{BatCmd, OhMyPoshCmd, TmuxCmd, ZshCmd},
}
