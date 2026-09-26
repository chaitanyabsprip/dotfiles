//go:build network

package install

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGhReleaseBareBinary(t *testing.T) {
	orig := BinDir
	BinDir = t.TempDir()
	defer func() { BinDir = orig }()

	err := GhRelease(`mikefarah/yq`, `yq`, func(_, goos, goarch string) string {
		return fmt.Sprintf(`yq_%s_%s`, goos, goarch)
	})
	if err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command(filepath.Join(BinDir, `yq`), `--version`).CombinedOutput(); err != nil {
		t.Fatalf("yq --version: %v: %s", err, out)
	}
}

func TestGhReleaseTarball(t *testing.T) {
	orig := BinDir
	BinDir = t.TempDir()
	defer func() { BinDir = orig }()

	arch := map[string]string{`amd64`: `x86_64`, `arm64`: `aarch64`}
	sys := map[string]string{`linux`: `unknown-linux-gnu`, `darwin`: `apple-darwin`}
	err := GhRelease(`sharkdp/fd`, `fd`, func(tag, goos, goarch string) string {
		return fmt.Sprintf(`fd-%s-%s-%s.tar.gz`, tag, arch[goarch], sys[goos])
	})
	if err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command(filepath.Join(BinDir, `fd`), `--version`).CombinedOutput(); err != nil {
		t.Fatalf("fd --version: %v: %s", err, out)
	}
}
