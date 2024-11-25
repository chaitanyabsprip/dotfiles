package tmux

import (
	"fmt"
	"os"
	"strings"

	"github.com/rwxrob/bonzai/fn/filt"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/to"
)

func IsActive() bool {
	return len(os.Getenv("TMUX")) > 0 ||
		len(run.Out(`pgrep`, `tmux`)) > 0
}

func ListSessions() []string {
	return to.Lines(
		run.Out(
			`tmux`,
			`ls`,
			`-F`,
			`#{session_name}=#{session_path}`,
		),
	)
}

type Session struct {
	Name string
	Path string
}

func NewSession(opts Session) error {
	fmt.Println("name:", opts.Name, "path:", opts.Path)
	return run.Exec(
		`tmux`,
		`new-session`,
		`-ds`,
		opts.Name,
		`-c`,
		opts.Path,
	)
}

func SessionExists(opts Session) bool {
	if len(opts.Name) > 0 {
		return len(filt.HasPrefix(ListSessions(), opts.Name)) > 0
	} else if len(opts.Path) > 0 {
		return len(filt.HasSuffix(ListSessions(), opts.Path)) > 0
	}
	return false
}

func FindSession(opts Session) (string, string) {
	var results []string
	if len(opts.Name) > 0 {
		results = filt.HasPrefix(ListSessions(), opts.Name)
	} else if len(opts.Path) > 0 {
		results = filt.HasPrefix(ListSessions(), opts.Path)
	}
	if len(results) == 0 {
		return "", ""
	}
	parts := strings.Split(results[0], "=")
	if len(parts) < 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func SwitchClient(name string) error {
	return run.Exec(`tmux`, `switch-client`, `-t`, name)
}

func RenameSession(old, new string) error {
	return run.Exec(`tmux`, `rename-session`, `-t`, old, new)
}
