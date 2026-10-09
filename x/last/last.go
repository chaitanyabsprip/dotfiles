// Package last finds and acts on the most recently modified file or
// directory in a directory, defaulting to $DOWNLOADS.
package last

import (
	"cmp"
	"os"
	"slices"
	"path/filepath"
)

// find returns the path to the most recently modified non-hidden entry
// directly inside dir for which match returns true (match nil means any
// entry). A matching subdirectory counts by its own mtime, not its
// newest file's. Returns "" if dir has no such entry.
func find(dir string, match func(isDir bool) bool) (string, error) {
	paths, err := findN(dir, 1, match)
	if err != nil || len(paths) == 0 {
		return ``, err
	}
	return paths[0], nil
}

// findN is find for up to n entries, newest first. Ties keep
// directory (name) order.
func findN(dir string, n int, match func(isDir bool) bool) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	type entry struct {
		name string
		mod  int64
	}
	var found []entry
	for _, e := range entries {
		if e.Name()[0] == '.' || (match != nil && !match(e.IsDir())) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		found = append(found, entry{e.Name(), info.ModTime().UnixNano()})
	}
	slices.SortStableFunc(found, func(a, b entry) int { return cmp.Compare(b.mod, a.mod) })
	paths := make([]string, 0, min(n, len(found)))
	for _, e := range found[:min(n, len(found))] {
		paths = append(paths, filepath.Join(dir, e.name))
	}
	return paths, nil
}

// Find returns the path to the newest non-hidden entry directly inside
// dir, regardless of whether it is a file or a directory.
func Find(dir string) (string, error) { return find(dir, nil) }

// FindFile returns the path to the newest non-hidden regular file (or
// other non-directory entry) directly inside dir.
func FindFile(dir string) (string, error) {
	return find(dir, isFileMatch)
}

// FindDir returns the path to the newest non-hidden subdirectory
// directly inside dir.
func FindDir(dir string) (string, error) {
	return find(dir, isDirMatch)
}

func isDirMatch(isDir bool) bool  { return isDir }
func isFileMatch(isDir bool) bool { return !isDir }
