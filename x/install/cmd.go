// Package install provides commands for installing and setting up various tools
// and applications used in the dotfiles ecosystem. It includes installers for
// terminal utilities like Bat, Oh My Posh, Tmux, and Zsh.
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
