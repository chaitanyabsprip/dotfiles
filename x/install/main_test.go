package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Chaitanyabsprip/dotfiles/internal/testenv"
)

// TestMain keeps installs out of the real ~/.local/bin. BinDir is derived
// from HOME at startup, so it is repointed here too.
func TestMain(m *testing.M) {
	home, cleanup, err := testenv.Isolate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	orig := BinDir
	BinDir = filepath.Join(home, `.local`, `bin`)
	code := m.Run()
	BinDir = orig
	cleanup()
	os.Exit(code)
}

func TestSandboxed(t *testing.T) {
	if !strings.Contains(os.Getenv(`HOME`), `dotfiles-test-home-`) {
		t.Errorf("not sandboxed: HOME=%q", os.Getenv(`HOME`))
	}
	if !strings.Contains(BinDir, `dotfiles-test-home-`) {
		t.Errorf("BinDir = %q points outside the sandbox", BinDir)
	}
}
