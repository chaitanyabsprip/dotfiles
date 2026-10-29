// Package tmux
package tmux

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/chaitanyabsprip/bonzai"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
	"github.com/Chaitanyabsprip/dotfiles/pkg/fzf"
	"github.com/Chaitanyabsprip/dotfiles/pkg/tmux"
	"github.com/Chaitanyabsprip/dotfiles/x/depends"
	"github.com/Chaitanyabsprip/dotfiles/x/work"
)

var SessionizerCmd = &bonzai.Cmd{
	Name:  `sessionizer`,
	Alias: `sz`,
	Short: `create a new or switch to an existing session`,
	Usage: `[path]`,
	Long: `
Switches to the tmux session for [path], or, without it, fuzzy-picks
a work directory with fzf. Creates the session (and resumes any
Claude Code sessions there) if none exists yet for that path.`,
	Init: func(x *bonzai.Cmd, args ...string) error {
		depends.On(nil, `fzf`, `tmux`)
		return nil
	},
	Do: func(x *bonzai.Cmd, args ...string) error {
		// should always exit with non-zero status
		if len(args) > 0 {
			Sessionizer(args[0])
			return nil
		}
		Sessionizer(``)
		return nil
	},
}

func Sessionizer(path string) error {
	newPath := path
	if len(newPath) == 0 {
		newPath = selectPath()
	}
	if len(newPath) == 0 {
		return fmt.Errorf(`no path selected`)
	}
	// tmux stores session paths absolute and clean; match that so "." or
	// "dir/" find their session instead of colliding with it.
	if abs, err := filepath.Abs(newPath); err == nil {
		newPath = abs
	}
	sessionName := resolveSessionName(newPath)

	if !tmux.IsActive() {
		if err := tmux.NewSession(
			tmux.Session{Name: sessionName, Path: newPath},
		); err != nil {
			return err
		}
		resumeClaude(sessionName)
		return nil
	}

	fopts := tmux.Session{Path: newPath}
	if oldName, _ := tmux.FindSession(fopts); len(oldName) > 0 {
		return tmux.SwitchClient(oldName)
	}
	if _, oldPath := tmux.FindSession(tmux.Session{Name: sessionName}); len(
		oldPath,
	) > 0 {
		collidingName := sessionName
		var oldSessionNewName string
		oldSessionNewName, sessionName = reconcileSessionName(oldPath, newPath)
		tmux.RenameSession(collidingName, oldSessionNewName)
	}
	err := tmux.NewSession(
		tmux.Session{Name: sessionName, Path: newPath},
	)
	if err != nil {
		return err
	}
	resumeClaude(sessionName)
	return tmux.SwitchClient(sessionName)
}

// resumeClaude runs the optional out-of-repo `claude-tmux` script, if
// installed, to resume Claude Code sessions remembered for this session.
func resumeClaude(sessionName string) {
	_ = exec.Command(`claude-tmux`, `restore`, sessionName).Run()
}

func resolveSessionName(path string) string {
	base := filepath.Base(path)
	if base == "root" {
		base = fmt.Sprintf("%s/r", filepath.Base(filepath.Dir(path)))
	}
	return strings.ReplaceAll(base, `.`, `_`)
}

func selectPath() string {
	dirs := work.CachedAllDirs(false)
	short := work.Shorten(dirs)
	out, err := fzf.Select(
		short,
		`--tmux`, `45%`,
		`--border`,
		`--border-label`, ` Sessionizer `,
		`--border-label-pos`, `6:bottom`,
	)
	if err != nil || len(out) == 0 {
		return ``
	}
	// Shorten keeps order, so map the pick back by index; a suffix match
	// would confuse `dotfiles` with `projects/.../dotfiles`.
	if i := slices.Index(short, out); i >= 0 {
		return dirs[i]
	}
	return ``
}

// reconcileSessionName returns the names for two sessions whose names
// collide: the old one is renamed to its shortest trailing path that
// differs from the new path's ("ss/main" vs "ui/main"), the new one
// keeps resolveSessionName. Compares sanitized names, so "foo.bar" and
// "foo_bar" don't both end up "foo_bar".
func reconcileSessionName(oldPath, newPath string) (string, string) {
	a, b := pathParts(oldPath), pathParts(newPath)
	for n := 1; n <= max(len(a), len(b)); n++ {
		if oldName := trailingName(a, n); oldName != trailingName(b, n) {
			return oldName, resolveSessionName(newPath)
		}
	}
	return `old-` + resolveSessionName(oldPath), resolveSessionName(newPath)
}

// pathParts splits path into its parts, with a path under $HOME rooted
// at "~" so names read "~/dotfiles", not "<user>/dotfiles".
func pathParts(path string) []string {
	if rest, ok := strings.CutPrefix(path, env.Home+`/`); ok && env.Home != `` {
		path = `~/` + rest
	}
	return strings.Split(strings.Trim(path, `/`), `/`)
}

// trailingName joins the last n path parts as a session name, with a
// trailing "root" shortened to "r" and dots made tmux-safe.
func trailingName(parts []string, n int) string {
	parts = slices.Clone(parts[max(0, len(parts)-n):])
	if last := len(parts) - 1; parts[last] == `root` {
		parts[last] = `r`
	}
	return strings.ReplaceAll(strings.Join(parts, `/`), `.`, `_`)
}
