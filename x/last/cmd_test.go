package last

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

func TestDownloadsDirPrefersEnvVar(t *testing.T) {
	env.Downloads = `/from/env`
	defer func() { env.Downloads = `` }()

	if got := downloadsDir(); got != `/from/env` {
		t.Errorf("downloadsDir() = %q, want /from/env", got)
	}
}

func TestDownloadsDirFallsBackToHome(t *testing.T) {
	env.Downloads = ``
	if got, want := downloadsDir(), filepath.Join(env.Home, `downloads`); got != want {
		t.Errorf("downloadsDir() = %q, want %q", got, want)
	}
}

// withDownloads points downloadsDir() at a fresh temp dir for one test.
func withDownloads(t *testing.T) string {
	t.Helper()
	downloads := filepath.Join(t.TempDir(), `downloads`)
	if err := os.MkdirAll(downloads, 0o755); err != nil {
		t.Fatal(err)
	}
	orig := env.Downloads
	env.Downloads = downloads
	t.Cleanup(func() { env.Downloads = orig })
	return downloads
}

func TestTransferMovesNewestFile(t *testing.T) {
	downloads := withDownloads(t)
	now := time.Now()
	touch(t, filepath.Join(downloads, `old.zip`), now.Add(-time.Hour))
	touch(t, filepath.Join(downloads, `report.pdf`), now)
	mkdirAt(t, filepath.Join(downloads, `newer-dir`), now.Add(time.Hour)) // skipped: FindFile only

	wd := t.TempDir()
	from, to, err := transfer(FindFile, wd, `renamed.pdf`, moveOp)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(downloads, `report.pdf`); from != want {
		t.Errorf("from = %q, want %q", from, want)
	}
	if want := filepath.Join(wd, `renamed.pdf`); to != want {
		t.Errorf("to = %q, want %q", to, want)
	}
	if _, err := os.Stat(to); err != nil {
		t.Errorf("moved file missing at %s: %v", to, err)
	}
	if _, err := os.Stat(from); !os.IsNotExist(err) {
		t.Errorf("original file still at %s: %v", from, err)
	}
}

func TestTransferNoMatch(t *testing.T) {
	withDownloads(t) // empty
	if _, _, err := transfer(FindFile, t.TempDir(), `x`, moveOp); err == nil {
		t.Error(`transfer with no matching entry should fail`)
	}
}

// Bare mv/cp (lookup = Find) must pick the newest entry of either type,
// not just files.
func TestTransferAnyTypeMovesNewestDir(t *testing.T) {
	downloads := withDownloads(t)
	now := time.Now()
	touch(t, filepath.Join(downloads, `older.pdf`), now.Add(-time.Hour))
	mkdirAt(t, filepath.Join(downloads, `newest-dir`), now)

	wd := t.TempDir()
	from, to, err := transfer(Find, wd, `moved-dir`, moveOp)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(downloads, `newest-dir`); from != want {
		t.Errorf("from = %q, want %q", from, want)
	}
	if _, err := os.Stat(to); err != nil {
		t.Errorf("moved dir missing at %s: %v", to, err)
	}
}

func TestCopyOpFileLeavesOriginal(t *testing.T) {
	src := filepath.Join(t.TempDir(), `src.txt`)
	if err := os.WriteFile(src, []byte(`hi`), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(t.TempDir(), `dst.txt`)

	if err := copyOp(src, dst); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(src); err != nil {
		t.Errorf("original file removed by copy: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil || string(got) != `hi` {
		t.Errorf("copy contents = %q, %v; want \"hi\", nil", got, err)
	}
}

func TestCopyOpDirRecursesAndLeavesOriginal(t *testing.T) {
	src := filepath.Join(t.TempDir(), `srcdir`)
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	inner := filepath.Join(src, `inner.txt`)
	if err := os.WriteFile(inner, []byte(`hi`), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(t.TempDir(), `dstdir`)

	if err := copyOp(src, dst); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(inner); err != nil {
		t.Errorf("original dir contents removed by copy: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dst, `inner.txt`))
	if err != nil || string(got) != `hi` {
		t.Errorf("copied dir contents = %q, %v; want \"hi\", nil", got, err)
	}
}
