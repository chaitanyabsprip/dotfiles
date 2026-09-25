package starship

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai/futil"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

// installStarship uses starship's official install script rather than a
// distro package manager: it's unavailable via apt/dnf (verified against
// Debian, Ubuntu and Fedora), so mirrors the ohmyposh install path instead.
func installStarship() error {
	if ok, _ := have.Executable(`starship`); ok {
		fmt.Println(`starship is already installed`)
		return nil
	}
	binDir := oscfg.BinDir()
	if !futil.Exists(binDir) {
		if err := futil.CreateDir(binDir); err != nil {
			return err
		}
	}
	curl := exec.Command(`curl`, `-sS`, `https://starship.rs/install.sh`)
	sh := exec.Command(`sh`, `-s`, `--`, `-y`, `-b`, binDir)
	sh.Stdin, _ = curl.StdoutPipe()
	sh.Stdout = os.Stdout
	sh.Stderr = os.Stderr
	if err := curl.Start(); err != nil {
		return err
	}
	return sh.Run()
}
