package embed

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/manifest"
)

//go:embed testdata
var testFs embed.FS

const (
	embeddedPath    = `testdata/sample.txt`
	embeddedContent = "embedded-content\n"
)

func sampleDirEntry(t *testing.T) fs.DirEntry {
	t.Helper()
	entries, err := fs.ReadDir(testFs, `testdata`)
	if err != nil {
		t.Fatalf(`ReadDir: %v`, err)
	}
	return entries[0]
}

func writeDest(t *testing.T, content string) string {
	t.Helper()
	dest := filepath.Join(t.TempDir(), `sample.txt`)
	if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
		t.Fatalf(`seed dest: %v`, err)
	}
	return dest
}

func TestCopyFirstDeploy(t *testing.T) {
	dest := filepath.Join(t.TempDir(), `sample.txt`)
	m := manifest.Manifest{}
	if err := copy(testFs, m, sampleDirEntry(t), embeddedPath, dest); err != nil {
		t.Fatalf(`copy: %v`, err)
	}
	assertContent(t, dest, embeddedContent)
	if m[dest] != manifest.Hash([]byte(embeddedContent)) {
		t.Fatalf(`manifest not recorded for first deploy`)
	}
}

func TestCopyUpstreamChangeAppliedAutomatically(t *testing.T) {
	dest := writeDest(t, `old-deployed-content`)
	m := manifest.Manifest{dest: manifest.Hash([]byte(`old-deployed-content`))}
	if err := copy(testFs, m, sampleDirEntry(t), embeddedPath, dest); err != nil {
		t.Fatalf(`copy: %v`, err)
	}
	assertContent(t, dest, embeddedContent)
	if m[dest] != manifest.Hash([]byte(embeddedContent)) {
		t.Fatalf(`manifest not updated after upstream change`)
	}
}

func TestCopySkipsUserDriftWithoutForce(t *testing.T) {
	t.Setenv(`DOT_FORCE`, ``)
	dest := writeDest(t, `hand-edited-content`)
	m := manifest.Manifest{dest: manifest.Hash([]byte(`what-was-last-deployed`))}
	if err := copy(testFs, m, sampleDirEntry(t), embeddedPath, dest); err != nil {
		t.Fatalf(`copy: %v`, err)
	}
	assertContent(t, dest, `hand-edited-content`)
	if m[dest] != manifest.Hash([]byte(`what-was-last-deployed`)) {
		t.Fatalf(`manifest must not change when drift is skipped`)
	}
}

func TestCopyForceOverwritesUserDrift(t *testing.T) {
	t.Setenv(`DOT_FORCE`, `1`)
	dest := writeDest(t, `hand-edited-content`)
	m := manifest.Manifest{dest: manifest.Hash([]byte(`what-was-last-deployed`))}
	if err := copy(testFs, m, sampleDirEntry(t), embeddedPath, dest); err != nil {
		t.Fatalf(`copy: %v`, err)
	}
	assertContent(t, dest, embeddedContent)
	if m[dest] != manifest.Hash([]byte(embeddedContent)) {
		t.Fatalf(`manifest not updated after forced overwrite`)
	}
}

func TestCopyUntrackedFileMatchingEmbeddedIsAdopted(t *testing.T) {
	dest := writeDest(t, embeddedContent)
	m := manifest.Manifest{} // no manifest entry — file exists outside dot
	if err := copy(testFs, m, sampleDirEntry(t), embeddedPath, dest); err != nil {
		t.Fatalf(`copy: %v`, err)
	}
	assertContent(t, dest, embeddedContent)
	if m[dest] != manifest.Hash([]byte(embeddedContent)) {
		t.Fatalf(`manifest not recorded for adopted file`)
	}
}

func TestCopyUntrackedFileDifferingFromEmbeddedIsTreatedAsDrift(t *testing.T) {
	t.Setenv(`DOT_FORCE`, ``)
	dest := writeDest(t, `pre-existing-unmanaged-content`)
	m := manifest.Manifest{} // no manifest entry
	if err := copy(testFs, m, sampleDirEntry(t), embeddedPath, dest); err != nil {
		t.Fatalf(`copy: %v`, err)
	}
	assertContent(t, dest, `pre-existing-unmanaged-content`)
	if _, tracked := m[dest]; tracked {
		t.Fatalf(`manifest must not record a file that was skipped as drift`)
	}
}

func assertContent(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf(`read %s: %v`, path, err)
	}
	if string(got) != want {
		t.Fatalf(`%s content = %q, want %q`, path, got, want)
	}
}
