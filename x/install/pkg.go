package install

import (
	"fmt"
	"strings"

	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/x/distro"
)

// Pkg installs name via the current distro's package manager.
//
// overrides maps a manager key ("pacman", "apt", "dnf", "brew") to the
// package name on that manager, for tools where it differs from name. A
// brew cask (as opposed to a formula) is given as "cask:<name>". A
// manager missing from overrides falls back to name; an explicit empty
// string marks the package unavailable on that manager.
//
// An unsupported OS, or a package explicitly marked unavailable, prints a
// manual-install hint and returns nil rather than failing — matching the
// existing zsh/tmux install behavior so one unsupported tool doesn't abort
// `dot install all`.
func Pkg(name string, overrides map[string]string) error {
	pkgName := func(mgr string) (string, bool) {
		if v, ok := overrides[mgr]; ok {
			return v, v != ``
		}
		return name, true
	}
	switch distro.Name() {
	case `Arch Linux`:
		if p, ok := pkgName(`pacman`); ok {
			return WithRoot(`pacman`, `-S`, `--noconfirm`, p)
		}
	case `Ubuntu`, `Debian GNU/Linux`:
		if p, ok := pkgName(`apt`); ok {
			return WithRoot(`apt-get`, `install`, `-y`, p)
		}
	case `Fedora Linux`:
		if p, ok := pkgName(`dnf`); ok {
			return run.Exec(`dnf`, `install`, p, `-y`)
		}
	case `Darwin`:
		if p, ok := pkgName(`brew`); ok {
			if cask, isCask := strings.CutPrefix(p, `cask:`); isCask {
				return run.Exec(`brew`, `install`, `--cask`, cask)
			}
			return run.Exec(`brew`, `install`, p)
		}
	default:
		fmt.Printf("%s: unsupported or unconfigured operating system, install manually\n", name)
		return nil
	}
	fmt.Printf("%s: not available on this OS via a package manager, install manually\n", name)
	return nil
}
