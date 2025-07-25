// Package bin provides functionality for managing user scripts and binaries.
// It handles the setup and configuration of executable files in the user's local
// bin directory, offering commands to install and edit these scripts.
package bin

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

//go:embed bin
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bin`,
	Short: `manage scripts`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, editCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup bin directory`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		binDir := filepath.Join(env.HOME, ".local")
		return e.SetupAll(embedFs, "bin", binDir, nil)
	},
}

var editCmd = &bonzai.Cmd{
	Name:   "edit",
	Short:  `edit bin configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		if err := edit.Files(oscfg.BinDir()); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run bin setup")
		return nil
	},
}
