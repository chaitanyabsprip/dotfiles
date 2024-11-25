package tmux

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/pkg/fzf"
	"github.com/Chaitanyabsprip/dotfiles/pkg/tmux"
)

var SessionManagerCmd = &bonzai.Cmd{
	Name:  `session-manager`,
	Alias: `sm`,
	Short: `list and switch existing sessions`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			args = append(args, ``)
		}
		return SessionManager(args[0])
	},
}

func SessionManager(path string) error {
	selected := path
	if len(selected) == 0 {
		selected = selectSession()
	}
	count := len(tmux.ListSessions())
	if count == 1 {
		return run.SysExec(
			`tmux`,
			`display-message`,
			`Only one session`,
		)
	}
	current, err := tmux.CurrentSession()
	if err != nil {
		return err
	}
	if current.Name != selected && len(selected) > 0 {
		return tmux.SwitchClient(selected)
	}
	return nil
}

func selectSession() string {
	exe := run.ExeName()
	out, err := fzf.Select(tmux.ListSessionsF(`#{session_name}`),
		`--tmux`, `80%,90%`,
		`--preview`, fmt.Sprintf(`%s tmux x p s {}`, exe),
		`--preview-window`, `top,85%`,
	)
	if err != nil {
		return ``
	}
	return out
}
