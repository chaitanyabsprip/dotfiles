package waybar

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

//go:embed waybar
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `waybar`,
	Short: `waybar is a utility to manage waybar configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, editCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup waybar`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "waybar", oscfg.ConfigDir(), nil)
	},
}

var editCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit waybar configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(
			oscfg.ConfigDir(),
			"waybar",
			"config.jsonc",
		)
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run waybar setup")
		return nil
	},
}
