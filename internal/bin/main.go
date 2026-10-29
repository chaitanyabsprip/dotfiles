// Package bin provides functionality for managing user scripts and binaries.
// It handles the setup and configuration of executable files in the user's local
// bin directory, offering commands to install and edit these scripts.
package bin

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/chaitanyabsprip/bonzai"
	"github.com/chaitanyabsprip/bonzai/comp"
	"github.com/chaitanyabsprip/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed bin
var embedFs embed.FS

var SetupCmd = &bonzai.Cmd{
	Name: `setup`,

	Short: `setup bin directory`,
	Comp:  comp.Opts,
	// ~/dotfiles/bin, which shell/env.sh puts on PATH ahead of ~/.local/bin;
	// ~/.local/bin stays an untracked dir for installer drops.
	Do: func(x *bonzai.Cmd, _ ...string) error {
		return e.SetupAll(embedFs, "bin", oscfg.DotfilesDir(), nil)
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   "edit",
	Short:  `edit bin configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		if err := edit.Files(filepath.Join(oscfg.DotfilesDir(), "bin")); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run bin setup")
		return nil
	},
}
