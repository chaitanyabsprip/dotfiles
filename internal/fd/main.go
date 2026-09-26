// Package fd installs fd.
package fd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rwxrob/bonzai"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Short: `install fd`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		if ok, _ := have.Executable(`fd`); ok {
			fmt.Println(`fd is already installed`)
			return nil
		}
		err := install.Pkg(`fd`, map[string]string{`apt`: `fd-find`, `dnf`: `fd-find`})
		if err != nil {
			return err
		}
		// Debian installs the binary as fdfind.
		fdfind, err := exec.LookPath(`fdfind`)
		if err != nil {
			return nil
		}
		if err := os.MkdirAll(install.BinDir, 0o755); err != nil {
			return err
		}
		return os.Symlink(fdfind, filepath.Join(install.BinDir, `fd`))
	},
}
