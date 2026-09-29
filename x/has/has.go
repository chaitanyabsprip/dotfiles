// Package has walks up from a directory through its parents looking
// for an entry (a file, directory, ...) with a given name, the way
// POSIX test(1) checks it.
package has

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrNotFound is returned when no ancestor directory has entry.
var ErrNotFound = errors.New(`not found in any parent directory`)

// Find walks up from dir through its parents (dir itself first)
// looking for entry, testing each candidate path the way typ names:
// 'e' exists, 'f' is a regular file, 'd' is a directory, 'L' is a
// symlink, 's' is a non-empty file. It returns the first match,
// closest to dir.
//
// ponytail: only the file-test types this codebase actually uses are
// supported (POSIX test(1) has ~20); add more if a real caller needs one.
func Find(dir, entry string, typ byte) (string, error) {
	test, err := testFor(typ)
	if err != nil {
		return ``, err
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return ``, err
	}
	for {
		path := filepath.Join(dir, entry)
		if test(path) {
			return path, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ``, ErrNotFound
		}
		dir = parent
	}
}

func testFor(typ byte) (func(path string) bool, error) {
	switch typ {
	case 'e':
		return func(path string) bool {
			_, err := os.Lstat(path)
			return err == nil
		}, nil
	case 'f':
		return func(path string) bool {
			info, err := os.Stat(path)
			return err == nil && info.Mode().IsRegular()
		}, nil
	case 'd':
		return func(path string) bool {
			info, err := os.Stat(path)
			return err == nil && info.IsDir()
		}, nil
	case 'L':
		return func(path string) bool {
			info, err := os.Lstat(path)
			return err == nil && info.Mode()&os.ModeSymlink != 0
		}, nil
	case 's':
		return func(path string) bool {
			info, err := os.Stat(path)
			return err == nil && info.Size() > 0
		}, nil
	default:
		return nil, fmt.Errorf(`unsupported file test type %q`, typ)
	}
}
