package tmux

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/tmux/icon"
)

var XCmd = &bonzai.Cmd{
	Name:  `tmux`,
	Alias: `x`,
	Short: `tmux utility commands`,
	Long: `
Utility commands for working with tmux: session and pane management,
the sessionizer, the harpoon session bookmarks, and the statusline
helpers (gitmux, icon). See 'tmux help' for the individual commands.`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{
		KillCmd,
		SessionizerCmd,
		PreviewCmd,
		SessionManagerCmd,
		NotesCmd,
		icon.Cmd,
		SuspendCmd,
		GitmuxCmd,
		HarpoonCmd,
		help.Cmd,
	},
	Do: func(x *bonzai.Cmd, _ ...string) error {
		fmt.Printf("%s - %s\n\n", x.Name, x.Short)
		fmt.Println(`COMMANDS:`)
		for _, c := range x.Cmds {
			fmt.Printf("  %-15s - %s\n", c.Name, c.Short)
		}
		return nil
	},
}
