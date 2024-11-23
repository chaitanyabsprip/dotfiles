package tmux

import (
	"embed"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
	"github.com/Chaitanyabsprip/dot/x/install"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
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
	Name:  `tmx`,
	Short: `tmx is a utility to manage tmux configuration and related scripts`,
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
	Short: `setup tmux copies configuration files to config directory`,
	Long:  ``,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		return e.SetupAll(embedFs, "tmux", oscfg.ConfigDir(), nil)
	},
}
