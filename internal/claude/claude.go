// Package claude wires Claude Code session tracking into tmux, so
// tmux-resurrect can resume the right `claude --resume <id>` into the right
// pane after a restore.
package claude

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `claude`,
	Short: `claude code tmux session tracking`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		HookCmd,
	},
}
