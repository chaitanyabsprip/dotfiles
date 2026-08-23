package workdirs

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

const defaultCacheTTL = 10 * time.Minute

// refreshWorkerEnv, when set to "1" in the environment, tells the
// current binary to act as a detached background cache-refresh worker
// instead of running its normal command tree. Entry points (cmd/dot,
// cmd/x) must check this before dispatching to their bonzai command
// tree — see RefreshCacheWorker.
const refreshWorkerEnv = `WORKDIRS_CACHE_REFRESH_WORKER`

type cacheFile struct {
	Timestamp int64    `json:"timestamp"`
	Dirs      []string `json:"dirs"`
}

func cachePath() string {
	return filepath.Join(oscfg.CacheDir(), `dot`, `workdirs.json`)
}

func refreshLockPath() string {
	return cachePath() + `.refreshing`
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

// CachedAllDirs returns AllDirs(), stale-while-revalidate: a cached
// result (of any age) is returned immediately, and a detached
// background process is kicked off to refresh the cache if it's older
// than the TTL (default 10m, override with WORKDIRS_CACHE_TTL in
// seconds) — the caller never blocks waiting for that refresh. Only
// the very first call ever (no cache on disk yet) blocks, since
// there's nothing to serve in the meantime.
//
// Caching is opt-in (see cacheEnabled); with it disabled this is
// equivalent to AllDirs(). Set refresh=true, or WORKDIRS_REFRESH=1 in
// the environment, to force a synchronous fresh recomputation instead.
func CachedAllDirs(refresh bool) []string {
	if !cacheEnabled() {
		return AllDirs()
	}
	if !refresh && len(os.Getenv(`WORKDIRS_REFRESH`)) > 0 {
		refresh = true
	}
	if refresh {
		dirs := AllDirs()
		writeCache(dirs)
		return dirs
	}

	dirs, exists, fresh := readCache()
	if !exists {
		dirs = AllDirs()
		writeCache(dirs)
		return dirs
	}
	if !fresh {
		spawnBackgroundRefresh()
	}
	return dirs
}

// RefreshCacheWorker recomputes AllDirs() and writes the cache. It is
// meant to run as a detached background process spawned by
// spawnBackgroundRefresh, not called directly in normal operation.
// Entry points should call this and return immediately when
// refreshWorkerEnv is set, before dispatching to their command tree.
func RefreshCacheWorker() {
	defer os.Remove(refreshLockPath())
	writeCache(AllDirs())
}

// IsRefreshWorker reports whether the current process was launched as
// a background cache-refresh worker (see RefreshCacheWorker).
func IsRefreshWorker() bool {
	return os.Getenv(refreshWorkerEnv) == `1`
}

// spawnBackgroundRefresh launches a detached copy of the current
// executable to run RefreshCacheWorker, then returns without waiting.
// A lock file guards against piling up redundant refreshers if called
// repeatedly while a refresh is already in flight; a lock older than
// its own generous max age is treated as abandoned (e.g. the worker
// crashed) and ignored.
func spawnBackgroundRefresh() {
	const maxLockAge = 2 * time.Minute

	lock := refreshLockPath()
	if info, err := os.Stat(lock); err == nil {
		if time.Since(info.ModTime()) < maxLockAge {
			return
		}
	}
	if err := os.MkdirAll(filepath.Dir(lock), 0o755); err != nil {
		return
	}
	if err := os.WriteFile(lock, nil, 0o644); err != nil {
		return
	}

	exe, err := os.Executable()
	if err != nil {
		os.Remove(lock)
		return
	}
	cmd := exec.Command(exe)
	cmd.Env = append(os.Environ(), refreshWorkerEnv+`=1`)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		os.Remove(lock)
	}
}

// readCache returns the cached dirs, whether a cache file exists at
// all, and whether it's still within the TTL. A missing or corrupt
// cache reports exists=false; a present-but-corrupt cache is treated
// as absent rather than erroring.
func readCache() (dirs []string, exists, fresh bool) {
	data, err := os.ReadFile(cachePath())
	if err != nil {
		return nil, false, false
	}
	var cf cacheFile
	if err := json.Unmarshal(data, &cf); err != nil {
		return nil, false, false
	}
	fresh = time.Since(time.Unix(cf.Timestamp, 0)) <= cacheTTL()
	return cf.Dirs, true, fresh
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
