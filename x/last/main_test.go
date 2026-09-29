package last

import (
	"fmt"
	"os"
	"testing"

	"github.com/Chaitanyabsprip/dotfiles/internal/testenv"
)

// TestMain keeps tests away from the user's real HOME and downloads.
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
