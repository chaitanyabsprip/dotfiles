package tmux

import (
	"embed"
	"fmt"
	"path"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/x/install"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

// TODO(me):
// - dependencies
// - live reload

//go:embed all:tmux
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `tmux`,
	Short: `manage tmux configuration and related scripts`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		initCmd,
		runCmd,
		install.TmuxCmd.WithName(`install`),
		editCmd,
	},
}

var editCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit tmux configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(oscfg.ConfigDir(), "tmux", "tmux.conf")
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run tmux init")
		return nil
	},
}

var initCmd = &bonzai.Cmd{
	Name:  `init`,
	Opts:  `slim|quik|full`,
	Short: `setup tmux to a specific level of configuration`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			args = append(args, `slim`)
		}
		mode := args[0]
		if mode == `slim` || mode == `quik` || mode == `full` {
			err := install.Tmux()
			if err != nil {
				return err
			}
			err = e.SetupAll(embedFs, "tmux", oscfg.ConfigDir(), nil)
			if err != nil {
				return err
			}
		}
		if mode == `quik` || mode == `full` {
			// install harpoon, pbc
		}
		if mode == `full` {
		}
		return nil
	},
}
