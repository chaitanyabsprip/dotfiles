package distro

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chaitanyabsprip/bonzai/futil"
)

// ID and Version read ID/VERSION_ID, not NAME ("Debian GNU/Linux"), and
// don't confuse ID with VERSION_ID or ID_LIKE.
func TestIDAndVersionFromOSRelease(t *testing.T) {
	f := filepath.Join(t.TempDir(), `os-release`)
	release := "PRETTY_NAME=\"Debian GNU/Linux 12 (bookworm)\"\n" +
		"NAME=\"Debian GNU/Linux\"\nVERSION_ID=\"12\"\nID=debian\nID_LIKE=\"ubuntu\"\n"
	if err := os.WriteFile(f, []byte(release), 0o644); err != nil {
		t.Fatal(err)
	}
	orig := osReleaseFile
	osReleaseFile = f
	defer func() { osReleaseFile = orig }()

	if got := ID(); got != `debian` {
		t.Errorf("ID() = %q, want debian", got)
	}
	if got := Version(); got != `12` {
		t.Errorf("Version() = %q, want 12", got)
	}
}

// Without /etc/os-release or lsb_release (macOS), ID falls back to uname.
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
	if got, want := ID(), strings.TrimSpace(string(out)); got != want {
		t.Errorf("ID() = %q, want %q", got, want)
	}
}
