package hypr

import (
	"embed"
	"fmt"
	"path"
	"runtime"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed hypr
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `hypr`,
	Short: `manage hypr configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, editCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Short: `setup hypr`,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		if runtime.GOOS != "linux" {
			return nil
		}
		return e.SetupAll(embedFs, `hypr`, oscfg.ConfigDir(), nil)
	},
}

var editCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit hypr configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		if runtime.GOOS != "linux" {
			fmt.Println("hypr is only supported on linux")
			return nil
		}
		filePath := path.Join(
			oscfg.ConfigDir(),
			"hypr",
			"hyprland.conf",
		)
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run hypr setup")
		return nil
	},
}
