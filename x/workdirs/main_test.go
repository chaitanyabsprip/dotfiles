package workdirs

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Chaitanyabsprip/dotfiles/internal/testenv"
)

var errGHInTests = errors.New(`gh called in a test without fakeGH`)

// TestMain runs the package with a throwaway HOME and gh seams that fail
// closed, so no test can touch the user's files, gh auth or the network.
func TestMain(m *testing.M) {
	_, cleanup, err := testenv.Isolate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ghPRView = func(string, string) (prInfo, error) { return prInfo{}, errGHInTests }
	ghPRCheckout = func(string, string) error { return errGHInTests }
	code := m.Run()
	cleanup()
	os.Exit(code)
}

// A test that forgets fakeGH must not reach the real gh or the network.
func TestGHFailsClosed(t *testing.T) {
	if _, err := ghPRView(t.TempDir(), `1`); !errors.Is(err, errGHInTests) {
		t.Errorf("ghPRView: err = %v, want errGHInTests", err)
	}
	if err := ghPRCheckout(t.TempDir(), `1`); !errors.Is(err, errGHInTests) {
		t.Errorf("ghPRCheckout: err = %v, want errGHInTests", err)
	}
}
