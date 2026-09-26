package dot

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Chaitanyabsprip/dotfiles/internal/testenv"
)

// TestMain keeps setup/install/edit tests away from the user's home.
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
	if !strings.Contains(os.Getenv(`HOME`), `dotfiles-test-home-`) {
		t.Errorf("not sandboxed: HOME=%q", os.Getenv(`HOME`))
	}
}
