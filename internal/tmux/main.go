package tmux

import (
	"embed"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/x/install"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

// TODO(me):
// - dependencies
// - live reload

// TODO(me):
// - Create options to control what level of setup is done
// - The options will be one of slim/quik/full
// - Slim: install only the files that are needed and skip the rest
// - Quik: install only the configuration files and skip the
//   optional dependencies
// - Full: install everything, all config files and all dependencies

//go:embed all:tmux
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `tmux`,
	Short: `manage tmux configuration and related scripts`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		setupCmd,
		runCmd,
		install.TmuxCmd.WithName(`install`),
	},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|full`,
	Short: `setup tmux to a specific level of configuration`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			args = append(args, `slim`)
		}
		mode := args[0]
		if mode == `slim` || mode == `quik` || mode == `full` {
			err := e.SetupAll(embedFs, "tmux", oscfg.ConfigDir(), nil)
			if err != nil {
				return err
			}
		}
		if mode == `quik` || mode == `full` {
			// install harpoon, pbc
		}
		if mode == `full` {
			// install gitmux
		}
		return e.SetupAll(embedFs, "tmux", oscfg.ConfigDir(), nil)
	},
}
