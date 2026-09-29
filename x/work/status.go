package work

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/arl/gitstatus"
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/term"
)

var statusCmd = &bonzai.Cmd{
	Name:   `status`,
	Short:  `show git status for every worktree of this project`,
	NoArgs: true,
	Long: `
Shows a one-line git status summary for every worktree of the project
containing the current directory (not arbitrary subdirectories — the
same worktrees ` + "`work trees`" + ` lists). A ` + "`pr/<number>`" + ` worktree also
shows that PR's gh state (OPEN/MERGED/CLOSED), or UNKNOWN if gh
couldn't be reached. Outside a git project, prints a warning instead
of failing.`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		err = walkStatuses(wd, func(s worktreeStatus) {
			branch := s.Branch
			if branch == `` {
				branch = `(detached)`
			}
			icon := colorize(term.Magenta, "\U000f062c")
			line := fmt.Sprintf("%s %-19s %-40s %s", icon, branch, Shorten([]string{s.Path})[0], summarize(s.Status))
			if s.PR != `` {
				line += fmt.Sprintf("  pr#%s %s", s.PR, prSymbol(s.PRState))
			}
			fmt.Println(line)
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, `warning: not inside a git project`)
		}
		return nil
	},
}

// worktreeStatus is one worktree's branch, path and git status. PR and
// PRState are set only when the worktree is a pr/<number> checkout: PR
// is the number, PRState its gh state (OPEN/MERGED/CLOSED), or
// "UNKNOWN" if gh couldn't be reached.
type worktreeStatus struct {
	Branch  string
	Path    string
	Status  *gitstatus.Status
	PR      string
	PRState string
}

// statuses returns the git status of every worktree in the project
// containing dir, skipping the bare admin worktree entry (if any).
func statuses(dir string) ([]worktreeStatus, error) {
	out := []worktreeStatus{}
	err := walkStatuses(dir, func(ws worktreeStatus) { out = append(out, ws) })
	return out, err
}

// walkStatuses computes the git status of every worktree in the
// project containing dir (skipping the bare admin worktree entry, if
// any) and calls fn as soon as each one is ready, rather than
// collecting them all first — each worktree involves its own gh PR
// lookup, so a project with several worktrees would otherwise sit
// silent until the slowest one finishes.
func walkStatuses(dir string, fn func(worktreeStatus)) error {
	p, err := openProject(dir)
	if err != nil {
		return err
	}
	wts, err := listWorktrees(p.Repo)
	if err != nil {
		return err
	}
	for _, wt := range wts {
		if wt.Bare {
			continue
		}
		st, err := statusOf(wt.Path)
		if err != nil {
			return err
		}
		ws := worktreeStatus{Branch: wt.Branch, Path: wt.Path, Status: st}
		if n := prFromPath(p.Root, wt.Path); n != `` {
			ws.PR = n
			ws.PRState = `UNKNOWN`
			if info, err := ghPRView(p.Repo, n); err == nil {
				ws.PRState = info.State
			}
		}
		fn(ws)
	}
	return nil
}

// statusOf runs gitstatus against dir. gitstatus always inspects the
// calling process's cwd, so this chdirs there and back.
func statusOf(dir string) (*gitstatus.Status, error) {
	pop, err := pushdir(dir)
	if err != nil {
		return nil, err
	}
	defer pop()
	return gitstatus.NewWithContext(context.Background())
}

// pushdir chdirs to dir and returns a func that chdirs back.
func pushdir(dir string) (popdir func() error, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if err := os.Chdir(dir); err != nil {
		return nil, err
	}
	return func() error { return os.Chdir(pwd) }, nil
}

// summarize renders a Status using the exact same component order,
// symbols and styles as internal/tmux's gitmux.conf — stats
// (insertions, deletions), then flags (clean, or staged, conflict,
// modified, stashed, untracked), then divergence (ahead, behind) —
// verified byte-for-byte against `x tmux gitmux`'s own output, each
// symbol followed by a space then its count, groups joined by a
// single space, so `work status` matches the tmux statusline exactly.
// Symbols are Go unicode escapes, not pasted glyphs: several are
// outside the BMP (nerd-font PUA codepoints) and silently collapse to
// nothing if pasted as literal raw-string runes.
func summarize(s *gitstatus.Status) string {
	var parts []string
	add := func(color, symbol string, n int) {
		if n > 0 {
			parts = append(parts, colorize(color, symbol+" "+fmt.Sprintf(`%d`, n)))
		}
	}
	add(term.Green, "", s.Insertions)
	add(term.Red, "", s.Deletions)
	if s.IsClean {
		add(term.White, "", s.NumStashed)
		parts = append(parts, colorize(term.Green, "✔"))
	} else {
		add(term.Green, "", s.NumStaged)
		add(term.Red, `!`, s.NumConflicts)
		add(term.Yellow, "", s.NumModified)
		add(term.White, "", s.NumStashed)
		add(term.White, `??`, s.NumUntracked)
	}
	add(term.Cyan, "\U000f0dbc", s.AheadCount)
	add(term.Magenta, "\U000f0db9", s.BehindCount)
	return strings.Join(parts, ` `)
}

// prSymbol renders a gh PR state as a single symbol: ✓ merged,
// ✗ closed, ● open, ? unknown.
func prSymbol(state string) string {
	switch state {
	case `MERGED`:
		return colorize(term.Green, `✓`)
	case `CLOSED`:
		return colorize(term.Red, `✗`)
	case `OPEN`:
		return colorize(term.Yellow, `●`)
	default:
		return `?`
	}
}

// colorize wraps s in color and Bold, reset after. Both collapse to ""
// when stdout isn't a terminal (term.IsInteractive), so piped output
// has no escape codes.
func colorize(color, s string) string {
	return term.Bold + color + s + term.Reset
}
