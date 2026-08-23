package claude

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

func settingsPath() string {
	return filepath.Join(env.Home, `.claude`, `settings.json`)
}

// mergeHook ensures settings["hooks"][event] contains one entry running
// command, without disturbing any other hooks already registered for that
// event (this file also carries hooks for other tools).
func mergeHook(settings map[string]any, event, command string) {
	hooksRaw, _ := settings[`hooks`].(map[string]any)
	if hooksRaw == nil {
		hooksRaw = map[string]any{}
		settings[`hooks`] = hooksRaw
	}
	entries, _ := hooksRaw[event].([]any)
	for _, entry := range entries {
		group, ok := entry.(map[string]any)
		if !ok {
			continue
		}
		inner, _ := group[`hooks`].([]any)
		for _, h := range inner {
			hookMap, ok := h.(map[string]any)
			if ok && hookMap[`command`] == command {
				return // already registered
			}
		}
	}
	entries = append(entries, map[string]any{
		`hooks`: []any{
			map[string]any{`type`: `command`, `command`: command},
		},
	})
	hooksRaw[event] = entries
}

// unmergeHook removes any entry running command from
// settings["hooks"][event], deleting the event key entirely if it ends up
// empty. Used to clean up hooks an earlier version of SetupCmd installed.
func unmergeHook(settings map[string]any, event, command string) {
	hooksRaw, _ := settings[`hooks`].(map[string]any)
	if hooksRaw == nil {
		return
	}
	entries, _ := hooksRaw[event].([]any)
	kept := entries[:0]
	for _, entry := range entries {
		group, ok := entry.(map[string]any)
		if !ok {
			kept = append(kept, entry)
			continue
		}
		inner, _ := group[`hooks`].([]any)
		remaining := inner[:0]
		for _, h := range inner {
			hookMap, ok := h.(map[string]any)
			if ok && hookMap[`command`] == command {
				continue
			}
			remaining = append(remaining, h)
		}
		if len(remaining) == 0 {
			continue
		}
		group[`hooks`] = remaining
		kept = append(kept, group)
	}
	if len(kept) == 0 {
		delete(hooksRaw, event)
	} else {
		hooksRaw[event] = kept
	}
}

// SetupCmd idempotently registers the SessionStart hook in
// ~/.claude/settings.json. It also removes any SessionEnd hook a previous
// version of this command installed there (see the comment on
// SessionStartCmd for why we no longer delete state on SessionEnd).
var SetupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Short: `register claude code session hooks`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		path := settingsPath()
		settings := map[string]any{}
		if data, err := os.ReadFile(path); err == nil {
			if err := json.Unmarshal(data, &settings); err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}

		mergeHook(settings, `SessionStart`, `dot claude hook start`)
		unmergeHook(settings, `SessionEnd`, `dot claude hook end`)

		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent(``, `  `)
		if err := enc.Encode(settings); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		return os.WriteFile(path, buf.Bytes(), 0o644)
	},
}
