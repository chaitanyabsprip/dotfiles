// Package manifest records the hash of every file dot has deployed, so
// `setup` can tell a user's local edit apart from an upstream config
// change instead of silently overwriting it. See VISION.md's "Drift
// Detection" section for the full design this implements.
package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

func path() string {
	return filepath.Join(oscfg.StateDir(), `dot`, `manifest.json`)
}

// Manifest maps a deployed file's absolute path to the sha256 hex digest
// recorded the last time dot wrote it.
type Manifest map[string]string

// Load reads the manifest from disk, returning an empty Manifest if none
// has been written yet.
func Load() (Manifest, error) { return load(path()) }

func load(p string) (Manifest, error) {
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return Manifest{}, nil
	}
	if err != nil {
		return nil, err
	}
	m := Manifest{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// Save writes the manifest to disk, creating its parent directory if
// needed.
func (m Manifest) Save() error { return m.save(path()) }

func (m Manifest) save(p string) error {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m, ``, `  `)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}

// Hash returns the sha256 hex digest of content.
func Hash(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}
