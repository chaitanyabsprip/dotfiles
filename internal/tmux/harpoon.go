package tmux

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

// bookmark is one harpoon entry, stored as "session=dir" or
// "session=dir:window.pane".
type bookmark struct {
	Session string
	Dir     string
	Target  string // "window.pane"; empty means the session's current pane
}

func (b bookmark) String() string {
	if b.Target == `` {
		return b.Session + `=` + b.Dir
	}
	return b.Session + `=` + b.Dir + `:` + b.Target
}

func parseBookmarks(data string) []bookmark {
	var out []bookmark
	for _, line := range strings.Split(data, "\n") {
		session, rest, ok := strings.Cut(strings.TrimSpace(line), `=`)
		if !ok || session == `` {
			continue
		}
		dir, target, _ := strings.Cut(rest, `:`)
		out = append(out, bookmark{Session: session, Dir: dir, Target: target})
	}
	return out
}

func formatBookmarks(bs []bookmark) string {
	var b strings.Builder
	for _, bm := range bs {
		b.WriteString(bm.String() + "\n")
	}
	return b.String()
}

// addBookmark appends b unless that exact bookmark is already tracked.
func addBookmark(bs []bookmark, b bookmark) []bookmark {
	if slices.Contains(bs, b) {
		return bs
	}
	return append(slices.Clone(bs), b)
}

// replaceBookmark puts b at 1-based slot, appending when slot is past the
// end, and drops any duplicate that creates.
func replaceBookmark(bs []bookmark, slot int, b bookmark) ([]bookmark, error) {
	if slot < 1 {
		return nil, fmt.Errorf(`slot must be 1 or more, got %d`, slot)
	}
	if slot > len(bs) {
		return addBookmark(bs, b), nil
	}
	out := slices.Clone(bs)
	out[slot-1] = b
	var deduped []bookmark
	for _, bm := range out {
		if !slices.Contains(deduped, bm) {
			deduped = append(deduped, bm)
		}
	}
	return deduped, nil
}

// removeBookmarks drops every bookmark for exactly the named session.
func removeBookmarks(bs []bookmark, session string) []bookmark {
	return slices.DeleteFunc(slices.Clone(bs), func(b bookmark) bool {
		return b.Session == session
	})
}

func harpoonFile() string {
	return harpoonFileFor(os.Getenv(`XDG_CACHE_HOME`), runtime.GOOS, env.Home)
}

// harpoonFileFor mirrors the standalone harpoon's cache location so both
// tools share one bookmarks file.
func harpoonFileFor(xdgCache, goos, home string) string {
	dir := filepath.Join(home, `.cache`)
	switch {
	case xdgCache != ``:
		dir = xdgCache
	case goos == `darwin`:
		dir = filepath.Join(home, `Library`, `Caches`)
	}
	return filepath.Join(dir, `.tmux-harpoon-sessions`)
}

func loadBookmarks(path string) ([]bookmark, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return parseBookmarks(string(data)), nil
}

func saveBookmarks(path string, bs []bookmark) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(formatBookmarks(bs)), 0o644)
}
