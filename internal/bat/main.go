// Package bat provides commands to manage the bat configuration.
package bat

import (
	"embed"
	"fmt"
	"path"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"
	"github.com/rwxrob/bonzai/run"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/pkg/with"
	"github.com/Chaitanyabsprip/dotfiles/x/install"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed bat
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bat`,
	Short: `manage bat configuration`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		setupCmd,
		install.BatCmd.WithName(`install`),
		editCmd,
	},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup bat`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) (err error) {
		err = e.SetupAll(embedFs, `bat`, oscfg.ConfigDir(), nil)
		if err != nil {
			return err
		}
		reset, err := with.Path(oscfg.BinDir())
		if err != nil {
			return err
		}
		defer func() { err = reset() }()
		return run.Exec(`bat`, `cache`, `--build`)
	},
}

var editCmd = &bonzai.Cmd{
	Name:   "edit",
	Short:  `edit bat configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		cfgDir := oscfg.ConfigDir()
		filePath := path.Join(cfgDir, "bat", "config")
		err := edit.Files(filePath)
		if err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run bat setup")
		return nil
	},
}
