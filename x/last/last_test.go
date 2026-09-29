package last

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// touch creates path with the given mtime.
func touch(t *testing.T, path string, mtime time.Time) {
	t.Helper()
	if err := os.WriteFile(path, []byte(`x`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mtime, mtime); err != nil {
		t.Fatal(err)
	}
}

func mkdirAt(t *testing.T, path string, mtime time.Time) {
	t.Helper()
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, mtime, mtime); err != nil {
		t.Fatal(err)
	}
}

func TestFindAnyType(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	touch(t, filepath.Join(dir, `old.txt`), now.Add(-time.Hour))
	mkdirAt(t, filepath.Join(dir, `newer-dir`), now) // newest overall, and it's a dir
	touch(t, filepath.Join(dir, `.hidden`), now.Add(time.Hour))

	got, err := Find(dir)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, `newer-dir`); got != want {
		t.Errorf("Find(%s) = %q, want %q", dir, got, want)
	}
}

func TestFindFile(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	touch(t, filepath.Join(dir, `old.txt`), now.Add(-time.Hour))
	touch(t, filepath.Join(dir, `new.txt`), now)
	mkdirAt(t, filepath.Join(dir, `newer-dir`), now.Add(time.Hour)) // newer, but a dir
	touch(t, filepath.Join(dir, `.hidden`), now.Add(2*time.Hour))   // newest, but hidden

	got, err := FindFile(dir)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, `new.txt`); got != want {
		t.Errorf("FindFile(%s) = %q, want %q", dir, got, want)
	}
}

func TestFindDir(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	mkdirAt(t, filepath.Join(dir, `old-dir`), now.Add(-time.Hour))
	mkdirAt(t, filepath.Join(dir, `new-dir`), now)
	touch(t, filepath.Join(dir, `newer.txt`), now.Add(time.Hour)) // newer, but a file

	got, err := FindDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(dir, `new-dir`); got != want {
		t.Errorf("FindDir(%s) = %q, want %q", dir, got, want)
	}
}

func TestFindDirIgnoresSubdirContents(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	sub := filepath.Join(dir, `sub`)
	mkdirAt(t, sub, now.Add(-time.Hour))
	// A file newer than sub, but inside it, must not surface as the answer.
	touch(t, filepath.Join(sub, `inner`), now.Add(time.Hour))

	got, err := FindDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if want := sub; got != want {
		t.Errorf("FindDir(%s) = %q, want %q (the subdir itself, not its contents)", dir, got, want)
	}
}

func TestFindNoMatch(t *testing.T) {
	dir := t.TempDir()
	touch(t, filepath.Join(dir, `onlyfile.txt`), time.Now())

	if got, err := FindDir(dir); err != nil || got != `` {
		t.Errorf("FindDir(dir with no subdirs) = %q, %v; want \"\", nil", got, err)
	}
}

func TestFindMissingDirErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), `nope`)
	if _, err := Find(missing); err == nil {
		t.Error(`Find on a missing directory should return an error`)
	}
	if _, err := FindFile(missing); err == nil {
		t.Error(`FindFile on a missing directory should return an error`)
	}
	if _, err := FindDir(missing); err == nil {
		t.Error(`FindDir on a missing directory should return an error`)
	}
}
