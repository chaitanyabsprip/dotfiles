// Package alacritty provides commands to manage the Alacritty terminal
// emulator configuration.
package alacritty

import (
	"embed"
	"fmt"
	"path"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed alacritty
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `alacritty`,
	Short: `manage alacritty configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, editCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup alacritty`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		return e.SetupAll(embedFs, "alacritty", oscfg.ConfigDir(), nil)
	},
}

var editCmd = &bonzai.Cmd{
	Name:   "edit",
	Short:  `edit alacritty configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		cfgDir := oscfg.ConfigDir()
		filePath := path.Join(cfgDir, "alacritty", "alacritty.toml")
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run alacritty setup")
		return nil
	},
}
