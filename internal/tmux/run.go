package tmux

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

var runCmd = &bonzai.Cmd{
	Name:  `run`,
	Alias: `x`,
	Short: `tmux x`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		KillCmd,
		SessionizerCmd,
		PreviewCmd,
		SessionManagerCmd,
		NotesCmd,
		IconNameCmd,
		SuspendCmd,
	},
}
