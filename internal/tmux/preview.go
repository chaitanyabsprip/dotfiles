package tmux

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/pkg/tmux"
)

var PreviewCmd = &bonzai.Cmd{
	Name:  `preview`,
	Alias: `p`,
	Short: `preview tmux session`,
	Cmds:  []*bonzai.Cmd{previewTreeCmd, previewSingleCmd},
	Comp:  comp.Cmds,
}

var previewTreeCmd = &bonzai.Cmd{
	Name:    `tree`,
	Alias:   `t`,
	Short:   `preview tmux session tree`,
	MaxArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			args = append(args, ``)
		}
		return previewTree(args[0])
	},
}

var previewSingleCmd = &bonzai.Cmd{
	Name:    `session`,
	Alias:   `s`,
	Short:   `preview tmux session content`,
	NumArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return previewSession(args[0])
	},
}

func previewTree(string) error {
	fmt.Println(`not yet implemented`)
	// IDS := tmux.ListSessionsF(`#{session_id}`)
	// for _, s := range IDS {
	// 	S := filt.HasPrefix(
	// 		tmux.ListSessionsF(
	// 			`#{session_id}#{session_name}: #{T:tree_mode_format}`,
	// 		),
	// 		s,
	// 	)
	// }
	return nil
}

func previewSession(session string) error {
	sessionId, err := tmux.SessionID(session)
	if err != nil {
		return err
	}
	fmt.Println(run.Out(`tmux`, `capture-pane`, `-ep`, `-t`, sessionId))
	return nil
}
