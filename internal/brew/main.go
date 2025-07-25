// Package brew provides functionality for managing Homebrew configuration.
// It offers commands to setup and edit Brewfile configurations, primarily
// targeted at macOS systems where Homebrew is commonly used.
package brew

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

//go:embed Brewfile
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `brew`,
	Short: `brew is a utility to manage brew configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, editCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup brew`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if runtime.GOOS != "darwin" {
			return nil
		}
		return e.SetupAll(embedFs, "brew", oscfg.ConfigDir(), nil)
	},
}

var editCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit brew configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		cfgDir := oscfg.ConfigDir()
		filePath := path.Join(cfgDir, "brew", "Brewfile")
		if err := edit.Files(filePath); err != nil {
			fmt.Println(err)
			return err
		}
		println("rebuild binary")
		println("re run brew setup")
		return nil
	},
}
