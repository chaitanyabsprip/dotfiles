package claude

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/run"
)

// hookPayload is the subset of Claude Code's SessionStart/SessionEnd hook
// JSON (delivered on stdin) that we care about.
type hookPayload struct {
	SessionID      string `json:"session_id"`
	Cwd            string `json:"cwd"`
	TranscriptPath string `json:"transcript_path"`
}

func readPayload() (hookPayload, error) {
	var p hookPayload
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return p, err
	}
	if err := json.Unmarshal(data, &p); err != nil {
		return p, err
	}
	return p, nil
}

// tmuxAddress resolves the pane Claude was launched in to its tmux
// structural address, if it was launched inside tmux at all.
func tmuxAddress() (session string, window, pane int, ok bool) {
	tmuxPane := os.Getenv(`TMUX_PANE`)
	if len(tmuxPane) == 0 {
		return ``, 0, 0, false
	}
	out := strings.TrimSpace(run.Out(
		`tmux`, `display-message`, `-p`, `-t`, tmuxPane,
		"#{session_name}\t#{window_index}\t#{pane_index}",
	))
	fields := strings.Split(out, "\t")
	if len(fields) != 3 {
		return ``, 0, 0, false
	}
	w, err := strconv.Atoi(fields[1])
	if err != nil {
		return ``, 0, 0, false
	}
	p, err := strconv.Atoi(fields[2])
	if err != nil {
		return ``, 0, 0, false
	}
	return fields[0], w, p, true
}

// SessionStartCmd is the only lifecycle hook we register: SessionEnd fires
// on *any* process exit (Ctrl-C, /exit, and — critically — the SIGHUP every
// pane gets when tmux kills the server), so deleting state there would wipe
// the very data post-restore-all needs moments later. Cleanup instead
// happens lazily in SaveSession, which prunes entries whose transcript file
// has actually disappeared.
var SessionStartCmd = &bonzai.Cmd{
	Name:  `start`,
	Short: `record session's tmux pane (SessionStart)`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		payload, err := readPayload()
		if err != nil {
			return err
		}
		session, window, pane, ok := tmuxAddress()
		entry := Entry{
			SessionID:      payload.SessionID,
			Cwd:            payload.Cwd,
			TranscriptPath: payload.TranscriptPath,
		}
		if ok {
			entry.TmuxSession = session
			entry.TmuxWindow = window
			entry.TmuxPane = pane
		}
		return SaveSession(entry)
	},
}

var RestoreCmd = &bonzai.Cmd{
	Name:  `restore`,
	Short: `resume remembered sessions into tmux panes`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		sessionFilter := ``
		if len(args) > 0 {
			sessionFilter = args[0]
		}
		return ResumeSessions(sessionFilter)
	},
}

var HookCmd = &bonzai.Cmd{
	Name:  `hook`,
	Short: `claude code / tmux-resurrect hook entry points`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		SessionStartCmd,
		RestoreCmd,
	},
}
