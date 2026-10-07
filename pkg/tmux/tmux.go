// Package tmux provides a Go interface for interacting with tmux sessions.
// It includes functions for creating, managing, and manipulating tmux sessions,
// as well as retrieving information about the current tmux environment.
package tmux

import (
	"fmt"
	"strings"

	"github.com/rwxrob/bonzai/fn/maps"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/to"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

func IsActive() bool {
	return len(env.Tmux) > 0 ||
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

// parseSessionLines splits "name=path" lines (as ListSessions
// produces) into Sessions, so callers can match the name or path
// fields exactly instead of as a raw substring of the whole line —
// a literal prefix/suffix check on "name=path" false-matches, e.g.
// querying for session name "ss" would otherwise match "ss_ui=...".
func parseSessionLines(lines []string) []Session {
	sessions := make([]Session, 0, len(lines))
	for _, line := range lines {
		name, path, ok := strings.Cut(line, `=`)
		if !ok {
			continue
		}
		sessions = append(sessions, Session{Name: name, Path: path})
	}
	return sessions
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
	name, _ := FindSession(opts)
	return len(name) > 0
}

func FindSession(opts Session) (string, string) {
	return findSession(parseSessionLines(ListSessions()), opts)
}

func findSession(sessions []Session, opts Session) (string, string) {
	for _, s := range sessions {
		if len(opts.Name) > 0 && s.Name == opts.Name {
			return s.Name, s.Path
		}
		if len(opts.Path) > 0 && s.Path == opts.Path {
			return s.Name, s.Path
		}
	}
	return ``, ``
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
