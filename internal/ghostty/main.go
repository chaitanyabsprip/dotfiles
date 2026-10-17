// Package ghostty provides functionality for managing the Ghostty
// terminal emulator configuration.
package ghostty

import (
	"embed"
	"fmt"
	"path"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

//go:embed ghostty
var embedFs embed.FS

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		if ok, _ := have.Executable(`ghostty`); ok {
			fmt.Println(`ghostty is already installed`)
			return nil
		}
		// Not in the apt/dnf repos; pacman and Homebrew (cask) have it.
		return install.Pkg(`ghostty`, map[string]string{
			`apt`: ``, `dnf`: ``, `brew`: `cask:ghostty`,
		})
	},
}

var SetupCmd = &bonzai.Cmd{
	Name: `setup`,

	Short: `setup ghostty`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "ghostty", oscfg.ConfigDir(), nil)
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit ghostty configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(oscfg.ConfigDir(), "ghostty", "config")
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run ghostty setup")
		return nil
	},
}
