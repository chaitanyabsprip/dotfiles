package distro

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/rwxrob/bonzai/futil"
)

// Without /etc/os-release or lsb_release (macOS), Name falls back to uname.
func TestNameFallsBackToUname(t *testing.T) {
	if futil.Exists(`/etc/os-release`) {
		t.Skip(`/etc/os-release present`)
	}
	if _, err := exec.LookPath(`lsb_release`); err == nil {
		t.Skip(`lsb_release present`)
	}
	out, err := exec.Command(`uname`, `-s`).Output()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := Name(), strings.TrimSpace(string(out)); got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}
