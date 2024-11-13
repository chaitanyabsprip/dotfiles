package distro

import (
	"strings"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dot/x/have"
)

// Name returns the name of the distro for UNIX-like systems.
// Possible values (not exhaustive):
//
// - Debian based distros:
//   - Debian GNU/Linux
//   - Ubuntu
//   - Linux Mint
//   - Pop!_OS
//   - Kali GNU/Linux
//
// - Red Hat based distros:
//   - Fedora Linux
//   - CentOS Linux
//   - Rocky Linux
//   - AlmaLinux
//
// - Arch based distros:
//   - Arch Linux
//   - Manjaro Linux
//   - EndeavourOS
//
// - SUSE based distros:
//   - openSUSE Leap
//   - openSUSE Tumbleweed
//
// - Other distros:
//   - Alpine Linux
func Name() string {
	osReleaseFile := `/etc/os-release`
	if futil.Exists(osReleaseFile) {
		line, err := futil.FindString(osReleaseFile, "(?m)^NAME=.*")
		if err != nil {
			return ``
		}
		parts := strings.Split(line, `=`)
		if len(parts) < 2 {
			return ``
		}
		return strings.Trim(parts[1], `"`)
	} else if ok, _ := have.Executable(`lsb_release`); ok {
		return strings.TrimSpace(run.Out(`lsb_release -si`))
	} else {
		return strings.TrimSpace(run.Out(`uname -s`))
	}
}

// Version returns the version of the distro
// TODO(chaitanya): Unimplemented
func Version() string { return "" }
