package tmux

import (
	"fmt"
	"os"
	"strings"

	"github.com/rwxrob/bonzai/fn/filt"
	"github.com/rwxrob/bonzai/fn/maps"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/to"
)

func IsActive() bool {
	return len(os.Getenv(`TMUX`)) > 0 ||
		len(run.Out(`pgrep`, `tmux`)) > 0
}

func ListSessionsF(format string) []string {
	return maps.TrimSpace(
		to.Lines(run.Out(`tmux`, `ls`, `-F`, format)),
	)
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

func CurrentSession() (Session, error) {
	out := run.Out(
		`tmux`,
		`display-message`,
		`-p`,
		`#{session_name}=#{session_path}`,
	)
	parts := strings.Split(out, `=`)
	if len(parts) < 2 {
		return Session{}, fmt.Errorf(`no session found`)
	}
	return Session{Name: parts[0], Path: parts[1]}, nil
}

func NewSession(opts Session) error {
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
		return len(
			filt.HasPrefix(ListSessions(), opts.Name),
		) > 0
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
		return ``, ``
	}
	parts := strings.Split(results[0], `=`)
	if len(parts) < 2 {
		return ``, ``
	}
	return parts[0], parts[1]
}

func SwitchClient(name string) error {
	return run.Exec(`tmux`, `switch-client`, `-t`, name)
}

func RenameSession(old, new string) error {
	return run.Exec(`tmux`, `rename-session`, `-t`, old, new)
}

func SessionID(name string) (string, error) {
	out := strings.TrimSpace(run.Out(
		`tmux`,
		`ls`,
		`-F`,
		`#{session_id}`,
		`-f`,
		fmt.Sprintf(`#{==:#{session_name},%s}`, name),
	))
	if len(out) == 0 {
		return ``, fmt.Errorf(`unknown session: %s`, name)
	}
	return out, nil
}

func GetOption(name, fallback string) string {
	out := strings.TrimSpace(
		run.Out(`tmux`, `show-option`, `-gqv`, name),
	)
	if len(out) == 0 {
		return fallback
	}
	return out
}
