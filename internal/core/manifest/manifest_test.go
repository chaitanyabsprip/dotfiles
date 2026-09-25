package manifest

import (
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), `manifest.json`)

	m, err := load(statePath)
	if err != nil {
		t.Fatalf(`load on missing file: %v`, err)
	}
	if len(m) != 0 {
		t.Fatalf(`expected empty manifest, got %v`, m)
	}

	m[`/home/user/.zshrc`] = Hash([]byte(`content`))
	if err := m.save(statePath); err != nil {
		t.Fatalf(`save: %v`, err)
	}

	reloaded, err := load(statePath)
	if err != nil {
		t.Fatalf(`load after save: %v`, err)
	}
	if reloaded[`/home/user/.zshrc`] != m[`/home/user/.zshrc`] {
		t.Fatalf(`round trip mismatch: got %q, want %q`, reloaded[`/home/user/.zshrc`], m[`/home/user/.zshrc`])
	}
}

func TestHashDeterministic(t *testing.T) {
	a := Hash([]byte(`same`))
	b := Hash([]byte(`same`))
	if a != b {
		t.Fatalf(`Hash not deterministic: %q != %q`, a, b)
	}
	if a == Hash([]byte(`different`)) {
		t.Fatalf(`Hash collided for different content`)
	}
}
