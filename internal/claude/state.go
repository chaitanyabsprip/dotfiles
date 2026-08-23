// Package claude tracks Claude Code sessions per tmux pane so they can be
// resumed automatically when tmux-resurrect recreates the layout.
package claude

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

// Entry is the persisted record of one Claude Code session, correlated to
// the tmux pane it was started in.
type Entry struct {
	SessionID      string    `json:"session_id"`
	Cwd            string    `json:"cwd"`
	TmuxSession    string    `json:"tmux_session,omitempty"`
	TmuxWindow     int       `json:"tmux_window,omitempty"`
	TmuxPane       int       `json:"tmux_pane,omitempty"`
	TranscriptPath string    `json:"transcript_path,omitempty"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func stateDir() string {
	return filepath.Join(oscfg.CacheDir(), "claude-tmux", "sessions")
}

func entryPath(sessionID string) string {
	return filepath.Join(stateDir(), sessionID+".json")
}

// SaveSession writes (or overwrites) the entry for entry.SessionID, and
// opportunistically prunes any other entries whose transcript no longer
// exists on disk.
func SaveSession(entry Entry) error {
	dir := stateDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	entry.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(entryPath(entry.SessionID), data, 0o644); err != nil {
		return err
	}
	pruneStale()
	return nil
}

// RemoveSession deletes the entry for sessionID, if any.
func RemoveSession(sessionID string) error {
	err := os.Remove(entryPath(sessionID))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func allEntries() []Entry {
	files, err := os.ReadDir(stateDir())
	if err != nil {
		return nil
	}
	entries := make([]Entry, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(stateDir(), f.Name()))
		if err != nil {
			continue
		}
		var e Entry
		if err := json.Unmarshal(data, &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries
}

// pruneGrace is how long an entry is left alone before its transcript file
// is required to exist. Claude Code doesn't necessarily create the
// transcript on disk at the exact instant SessionStart fires, so pruning
// immediately would delete the entry SaveSession just wrote.
const pruneGrace = 5 * time.Minute

// pruneStale removes entries older than pruneGrace whose transcript file
// has disappeared, since that means the session is gone and will never be
// resumable again.
func pruneStale() {
	for _, e := range allEntries() {
		if len(e.TranscriptPath) == 0 || time.Since(e.UpdatedAt) < pruneGrace {
			continue
		}
		if _, err := os.Stat(e.TranscriptPath); os.IsNotExist(err) {
			RemoveSession(e.SessionID)
		}
	}
}

// EntriesForSession returns every entry saved for a tmux session, most
// recently updated first. Entries that duplicate a (cwd, relative window
// order) slot are expected to be deduplicated by the caller — see
// resume.go, which needs relative order rather than absolute window
// indices, since tmux renumbers windows on restore (`renumber-windows`).
func EntriesForSession(session string) []Entry {
	var matches []Entry
	for _, e := range allEntries() {
		if e.TmuxSession == session {
			matches = append(matches, e)
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].UpdatedAt.After(matches[j].UpdatedAt)
	})
	return matches
}
