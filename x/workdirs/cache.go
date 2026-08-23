package workdirs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

const defaultCacheTTL = 10 * time.Minute

type cacheFile struct {
	Timestamp int64    `json:"timestamp"`
	Dirs      []string `json:"dirs"`
}

func cachePath() string {
	return filepath.Join(oscfg.CacheDir(), `dot`, `workdirs.json`)
}

func cacheTTL() time.Duration {
	if raw := os.Getenv(`WORKDIRS_CACHE_TTL`); len(raw) > 0 {
		if secs, err := strconv.Atoi(raw); err == nil && secs >= 0 {
			return time.Duration(secs) * time.Second
		}
	}
	return defaultCacheTTL
}

// cacheEnabled reports whether on-disk caching is turned on. Caching
// is opt-in and off by default: walking the filesystem fresh every
// time is correct everywhere, but on some machines (e.g. ones running
// endpoint-security agents that tax every filesystem syscall) it's
// slow enough to be worth trading a little staleness for speed. Opt
// in per-machine with WORKDIRS_CACHE=1 in the environment.
func cacheEnabled() bool {
	switch os.Getenv(`WORKDIRS_CACHE`) {
	case `1`, `true`, `yes`:
		return true
	default:
		return false
	}
}

// AllDirs computes the full, deduplicated, sorted list of workdirs and
// worktrees by walking the filesystem. This is the expensive path;
// prefer CachedAllDirs for interactive use.
func AllDirs() []string {
	dirs := Workdirs()
	dirs = append(dirs, Worktrees()...)
	dirs = dedupe(dirs)
	sort.Strings(dirs)
	return dirs
}

// CachedAllDirs returns AllDirs(), served from an on-disk cache when
// caching is enabled (WORKDIRS_CACHE=1, off by default — see
// cacheEnabled) and the cache is younger than the TTL (default 10m,
// override with WORKDIRS_CACHE_TTL in seconds). Set refresh=true, or
// WORKDIRS_REFRESH=1 in the environment, to force a fresh
// recomputation regardless of cache age. With caching disabled this is
// equivalent to AllDirs().
func CachedAllDirs(refresh bool) []string {
	if !cacheEnabled() {
		return AllDirs()
	}
	if !refresh && len(os.Getenv(`WORKDIRS_REFRESH`)) > 0 {
		refresh = true
	}
	if !refresh {
		if dirs, ok := readCache(); ok {
			return dirs
		}
	}
	dirs := AllDirs()
	writeCache(dirs)
	return dirs
}

func readCache() ([]string, bool) {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return nil, false
	}
	var cf cacheFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, false
	}
	if time.Since(time.Unix(cf.Timestamp, 0)) > cacheTTL() {
		return nil, false
	}
	return cf.Dirs, true
}

func writeCache(dirs []string) {
	path := cachePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return
	}
	data, err := json.Marshal(cacheFile{
		Timestamp: time.Now().Unix(),
		Dirs:      dirs,
	})
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0o644)
}
