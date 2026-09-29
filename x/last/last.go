// Package last finds and acts on the most recently modified file or
// directory in a directory, defaulting to $DOWNLOADS.
package last

import (
	"os"
	"path/filepath"
)

// find returns the path to the most recently modified non-hidden entry
// directly inside dir for which match returns true (match nil means any
// entry). A matching subdirectory counts by its own mtime, not its
// newest file's. Returns "" if dir has no such entry.
func find(dir string, match func(isDir bool) bool) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ``, err
	}
	var newest string
	var newestMod int64
	for _, e := range entries {
		if e.Name()[0] == '.' || (match != nil && !match(e.IsDir())) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if mod := info.ModTime().UnixNano(); newest == `` || mod > newestMod {
			newest, newestMod = e.Name(), mod
		}
	}
	if newest == `` {
		return ``, nil
	}
	return filepath.Join(dir, newest), nil
}

// Find returns the path to the newest non-hidden entry directly inside
// dir, regardless of whether it is a file or a directory.
func Find(dir string) (string, error) { return find(dir, nil) }

// FindFile returns the path to the newest non-hidden regular file (or
// other non-directory entry) directly inside dir.
func FindFile(dir string) (string, error) {
	return find(dir, func(isDir bool) bool { return !isDir })
}

// FindDir returns the path to the newest non-hidden subdirectory
// directly inside dir.
func FindDir(dir string) (string, error) {
	return find(dir, func(isDir bool) bool { return isDir })
}
