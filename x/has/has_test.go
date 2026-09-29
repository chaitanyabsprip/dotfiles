package has

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindWalksUpToNearestMatch(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, `marker`))
	sub := filepath.Join(root, `a`, `b`, `c`)
	mkdir(t, sub)

	got, err := Find(sub, `marker`, 'e')
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, `marker`)
	if got != want {
		t.Errorf("Find() = %q, want %q", got, want)
	}
}

func TestFindPrefersClosestAncestor(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, `marker`))
	mid := filepath.Join(root, `a`)
	mkdir(t, filepath.Join(mid, `marker`))
	sub := filepath.Join(mid, `b`)
	mkdir(t, sub)

	got, err := Find(sub, `marker`, 'e')
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(mid, `marker`)
	if got != want {
		t.Errorf("Find() = %q, want closest match %q", got, want)
	}
}

func TestFindNotFound(t *testing.T) {
	sub := filepath.Join(t.TempDir(), `a`, `b`)
	mkdir(t, sub)
	if _, err := Find(sub, `nope-nowhere`, 'e'); err != ErrNotFound {
		t.Errorf("Find() error = %v, want ErrNotFound", err)
	}
}

func TestFindTypeDir(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, `thing`))
	sub := filepath.Join(root, `a`)
	mkdir(t, sub)
	writeFile(t, filepath.Join(sub, `thing`)) // a file, not a dir

	got, err := Find(sub, `thing`, 'd')
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, `thing`)
	if got != want {
		t.Errorf("Find() = %q, want the directory at %q, not the closer file", got, want)
	}
}

func TestFindTypeFile(t *testing.T) {
	root := t.TempDir()
	mkdir(t, filepath.Join(root, `thing`)) // a dir, not a file
	sub := filepath.Join(root, `a`)
	mkdir(t, sub)

	if _, err := Find(sub, `thing`, 'f'); err != ErrNotFound {
		t.Errorf("Find() error = %v, want ErrNotFound (dir shouldn't match 'f')", err)
	}
}

func TestFindInvalidType(t *testing.T) {
	if _, err := Find(t.TempDir(), `x`, 'z'); err == nil {
		t.Error(`Find() with an unsupported type should error`)
	}
}

func mkdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(`x`), 0o644); err != nil {
		t.Fatal(err)
	}
}
