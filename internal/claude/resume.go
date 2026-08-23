package claude

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/rwxrob/bonzai/run"
)

const paneFormat = "#{session_name}\t#{window_index}\t#{pane_index}\t#{pane_current_path}\t#{pane_current_command}"

type pane struct {
	session string
	window  int
	index   int
	path    string
	command string
}

func listPanes(sessionFilter string) []pane {
	args := []string{`tmux`, `list-panes`}
	if len(sessionFilter) > 0 {
		// -s: every pane in the session, not just the current window's.
		args = append(args, `-s`, `-t`, sessionFilter)
	} else {
		args = append(args, `-a`)
	}
	args = append(args, `-F`, paneFormat)
	out := run.Out(args...)

	var panes []pane
	for line := range strings.SplitSeq(out, "\n") {
		line = strings.TrimRight(line, "\r")
		if len(line) == 0 {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			continue
		}
		window, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		index, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}
		panes = append(panes, pane{
			session: fields[0],
			window:  window,
			index:   index,
			path:    fields[3],
			command: fields[4],
		})
	}
	return panes
}

type paneGroupKey struct {
	session string
	cwd     string
}

// ResumeSessions walks the tmux panes (all of them, or only sessionFilter's
// when non-empty) and, for every pane sitting at a plain shell with a
// remembered Claude session for its (tmux session, cwd), sends
// `claude --resume <id>` into it.
//
// Panes are matched by relative order within a (session, cwd) group, not by
// absolute window/pane index: `renumber-windows` (and tmux-resurrect's
// restore process in general) can shift a window's index across a restore,
// so exact-address matching is unreliable — but it always preserves
// relative order, since renumbering only compacts gaps. So the first saved
// session for a given project (by its original window/pane order) goes
// back into whichever pane ends up first among that project's restored
// panes, and so on. A saved session is never resumed into more than one
// pane.
func ResumeSessions(sessionFilter string) error {
	groups := map[paneGroupKey][]pane{}
	for _, p := range listPanes(sessionFilter) {
		if p.command == `claude` {
			continue
		}
		key := paneGroupKey{session: p.session, cwd: filepath.Clean(p.path)}
		groups[key] = append(groups[key], p)
	}

	for key, panes := range groups {
		sort.Slice(panes, func(i, j int) bool {
			if panes[i].window != panes[j].window {
				return panes[i].window < panes[j].window
			}
			return panes[i].index < panes[j].index
		})
		entries := entriesForGroup(key)
		for i := 0; i < len(panes) && i < len(entries); i++ {
			p, entry := panes[i], entries[i]
			target := p.session + ":" + strconv.Itoa(p.window) + "." + strconv.Itoa(p.index)
			if err := run.Exec(`tmux`, `send-keys`, `-t`, target, `claude --resume `+entry.SessionID, `Enter`); err != nil {
				return err
			}
		}
	}
	return nil
}

// entriesForGroup returns the saved entries for key.session whose cwd
// matches key.cwd, ordered by their originally-saved window/pane index —
// the same relative order the caller sorts live panes by.
func entriesForGroup(key paneGroupKey) []Entry {
	var matches []Entry
	for _, e := range EntriesForSession(key.session) {
		if filepath.Clean(e.Cwd) == key.cwd {
			matches = append(matches, e)
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].TmuxWindow != matches[j].TmuxWindow {
			return matches[i].TmuxWindow < matches[j].TmuxWindow
		}
		return matches[i].TmuxPane < matches[j].TmuxPane
	})
	return matches
}
