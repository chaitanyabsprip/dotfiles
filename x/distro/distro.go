package distro

import (
	"strings"

	"github.com/chaitanyabsprip/bonzai/futil"
	"github.com/chaitanyabsprip/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

// ID returns the distro's os-release ID for UNIX-like systems: a stable
// lowercase identifier, unlike NAME ("Debian GNU/Linux", "Fedora Linux")
// which varies in spelling. Possible values (not exhaustive):
//
//   - Debian based: debian, ubuntu, linuxmint, pop, kali
//   - Red Hat based: fedora, centos, rocky, almalinux
//   - Arch based: arch, manjaro, endeavouros
//   - SUSE based: opensuse-leap, opensuse-tumbleweed
//   - Other: alpine
//
// Falls back to lowercased `lsb_release -si`, then `uname -s` ("Darwin"
// on macOS).
func ID() string {
	return lookup(`ID`, `-si`, `-s`, strings.ToLower)
}

// Version returns the distro's os-release VERSION_ID ("24.04", "40"),
// falling back to `lsb_release -sr`, then `uname -r`. Rolling releases
// like Arch have no VERSION_ID and return "".
func Version() string {
	return lookup(`VERSION_ID`, `-sr`, `-r`, func(s string) string { return s })
}

// osReleaseFile is a var so tests can point it at a fixture.
var osReleaseFile = `/etc/os-release`

// lookup reads key from osReleaseFile, else runs lsb_release with
// lsbFlag (its output passed through lsbFix), else uname with unameFlag.
func lookup(key, lsbFlag, unameFlag string, lsbFix func(string) string) string {
	if futil.Exists(osReleaseFile) {
		line, err := futil.FindString(osReleaseFile, `(?m)^`+key+`=.*`)
		if err != nil {
			return ``
		}
		_, val, _ := strings.Cut(line, `=`)
		return strings.Trim(val, `"'`)
	}
	if ok, _ := have.Executable(`lsb_release`); ok {
		return lsbFix(strings.TrimSpace(run.Out(`lsb_release`, lsbFlag)))
	}
	return strings.TrimSpace(run.Out(`uname`, unameFlag))
}
