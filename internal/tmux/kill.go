package tmux

import (
	"fmt"
	"os"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/run"
)

func numOfSession() int {
	out := run.Out(`tmux`, `list-sessions`)
	return len(strings.Split(out, "\n")) - 1
}

func switchSession() error {
	count := numOfSession()
	if count < 1 {
		return fmt.Errorf("no tmux session found")
	}
	if count > 1 {
		run.Out(`tmux`, `switch-client -l`)
	} else {
		run.Out(`tmux`, `new-session`, `-s`, `home`, `-c`, os.Getenv(`HOME`))
	}
	return nil
}

func KillSession() error {
	currSession := run.Out(
		`tmux`,
		`display-message`,
		`-p`,
		`#{client_session}`,
	)
	err := switchSession()
	if err != nil {
		return err
	}
	return run.SysExec(
		`tmux`,
		`kill-session`,
		`-t`,
		strings.TrimSpace(currSession),
	)
}

var KillCmd = &bonzai.Cmd{
	Name:  `kill`,
	Alias: `k`,
	Short: `kill current tmux session`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return KillSession()
	},
}
