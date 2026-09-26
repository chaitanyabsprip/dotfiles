package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
	"github.com/Chaitanyabsprip/dotfiles/pkg/fzf"
	"github.com/Chaitanyabsprip/dotfiles/pkg/tmux"
)

var HarpoonCmd = &bonzai.Cmd{
	Name:  `harpoon`,
	Alias: `hp`,
	Short: `jump between bookmarked tmux sessions`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		hpAddCmd, hpPaneCmd, hpSetCmd, hpSetPaneCmd,
		hpRmCmd, hpListCmd, hpGoCmd, hpEditCmd,
	},
}

var hpAddCmd = &bonzai.Cmd{
	Name: `add`, Alias: `a`, Short: `bookmark the current session`, NoArgs: true,
	Do: func(_ *bonzai.Cmd, _ ...string) error { return track(false, 0) },
}

var hpPaneCmd = &bonzai.Cmd{
	Name: `pane`, Alias: `p`, Short: `bookmark the current pane`, NoArgs: true,
	Do: func(_ *bonzai.Cmd, _ ...string) error { return track(true, 0) },
}

var hpSetCmd = &bonzai.Cmd{
	Name: `set`, Alias: `r`, Short: `bookmark the current session at slot n`,
	Usage: `set <n>`, NumArgs: 1,
	Do: func(_ *bonzai.Cmd, args ...string) error { return trackAt(false, args[0]) },
}

var hpSetPaneCmd = &bonzai.Cmd{
	Name: `setpane`, Alias: `sp`, Short: `bookmark the current pane at slot n`,
	Usage: `setpane <n>`, NumArgs: 1,
	Do: func(_ *bonzai.Cmd, args ...string) error { return trackAt(true, args[0]) },
}

var hpRmCmd = &bonzai.Cmd{
	Name: `rm`, Alias: `d`, Short: `stop tracking a session (default: current)`,
	Usage: `rm [session]`, MaxArgs: 1,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		session := ``
		if len(args) == 1 {
			session = args[0]
		}
		return untrack(session)
	},
}

var hpListCmd = &bonzai.Cmd{
	Name: `list`, Alias: `l`, Short: `pick a bookmark with fzf`, NoArgs: true,
	Do: func(_ *bonzai.Cmd, _ ...string) error { return pick() },
}

var hpGoCmd = &bonzai.Cmd{
	Name: `go`, Alias: `s`, Short: `switch to bookmark n`, Usage: `go <n>`, NumArgs: 1,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		n, err := slot(args[0])
		if err != nil {
			return err
		}
		return goTo(n)
	},
}

var hpEditCmd = &bonzai.Cmd{
	Name: `edit`, Alias: `e`, Short: `edit the bookmarks file`, NoArgs: true,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		editor := os.Getenv(`EDITOR`)
		if editor == `` {
			editor = `vi`
		}
		return run.Exec(`tmux`, `display-popup`, `-E`, editor+` "`+harpoonFile()+`"`)
	},
}

func slot(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return 0, fmt.Errorf(`slot must be a number from 1, got %q`, s)
	}
	return n, nil
}

// updateBookmarks applies fn to the saved bookmarks and saves the result.
func updateBookmarks(fn func([]bookmark) ([]bookmark, error)) error {
	path := harpoonFile()
	bs, err := loadBookmarks(path)
	if err != nil {
		return err
	}
	if bs, err = fn(bs); err != nil {
		return err
	}
	return saveBookmarks(path, bs)
}

// currentBookmark describes the client's session, and pane if withPane.
func currentBookmark(withPane bool) (bookmark, error) {
	format := `#{session_name}=#{session_path}`
	if withPane {
		format += `:#{window_index}.#{pane_index}`
	}
	bs := parseBookmarks(run.Out(`tmux`, `display-message`, `-p`, format))
	if len(bs) != 1 {
		return bookmark{}, fmt.Errorf(`not inside tmux`)
	}
	return bs[0], nil
}

func track(withPane bool, at int) error {
	b, err := currentBookmark(withPane)
	if err != nil {
		return err
	}
	err = updateBookmarks(func(bs []bookmark) ([]bookmark, error) {
		if at == 0 {
			return addBookmark(bs, b), nil
		}
		return replaceBookmark(bs, at, b)
	})
	if err != nil {
		return err
	}
	return run.Exec(`tmux`, `display-message`, `Tracking `+b.Session)
}

func trackAt(withPane bool, s string) error {
	n, err := slot(s)
	if err != nil {
		return err
	}
	return track(withPane, n)
}

func untrack(session string) error {
	if session == `` {
		b, err := currentBookmark(false)
		if err != nil {
			return err
		}
		session = b.Session
	}
	err := updateBookmarks(func(bs []bookmark) ([]bookmark, error) {
		return removeBookmarks(bs, session), nil
	})
	if err != nil {
		return err
	}
	return run.Exec(`tmux`, `display-message`, `Stopped tracking `+session)
}

// goTo switches to bookmark n, recreating its session if it is gone.
func goTo(n int) error {
	path := harpoonFile()
	bs, err := loadBookmarks(path)
	if err != nil {
		return err
	}
	if n > len(bs) {
		return run.Exec(`tmux`, `display-message`, fmt.Sprintf(`No bookmark %d`, n))
	}
	b := bs[n-1]
	if exec.Command(`tmux`, `has-session`, `-t`, `=`+b.Session).Run() != nil {
		if err := tmux.NewSession(tmux.Session{Name: b.Session, Path: b.Dir}); err != nil {
			return err
		}
		if b.Target != `` { // the saved pane cannot exist in a new session
			b.Target = ``
			bs[n-1] = b
			if err := saveBookmarks(path, bs); err != nil {
				return err
			}
		}
	}
	target := `=` + b.Session
	if b.Target != `` {
		target += `:` + b.Target
	}
	return tmux.SwitchClient(target)
}

func pick() error {
	bs, err := loadBookmarks(harpoonFile())
	if err != nil {
		return err
	}
	if len(bs) == 0 {
		return run.Exec(`tmux`, `display-message`, `No bookmarks yet`)
	}
	lines := make([]string, len(bs))
	for i, b := range bs {
		dir := strings.Replace(b.Dir, env.Home, `~`, 1)
		lines[i] = fmt.Sprintf(`%d  %s  %s`, i+1, b.Session, dir)
	}
	out, err := fzf.Select(lines, `--tmux`, `50%,50%`, `--border-label`, ` Harpoon `)
	if err != nil || out == `` {
		return nil // cancelled
	}
	n, err := strconv.Atoi(strings.Fields(out)[0])
	if err != nil {
		return err
	}
	return goTo(n)
}
