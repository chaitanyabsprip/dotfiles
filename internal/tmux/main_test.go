package tmux

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Chaitanyabsprip/dotfiles/internal/testenv"
)

// TestMain keeps tests away from the user's home, harpoon file and tmux
// server.
func TestMain(m *testing.M) {
	_, cleanup, err := testenv.Isolate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestSandboxed(t *testing.T) {
	if !strings.Contains(os.Getenv(`HOME`), `dotfiles-test-home-`) || os.Getenv(`TMUX`) != `` {
		t.Errorf("not sandboxed: HOME=%q TMUX=%q", os.Getenv(`HOME`), os.Getenv(`TMUX`))
	}
	if !strings.Contains(harpoonFile(), `dotfiles-test-home-`) {
		t.Errorf("harpoonFile() = %q points outside the sandbox", harpoonFile())
	}
}
