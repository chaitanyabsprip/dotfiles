package tmux

import (
	bonzai "github.com/rwxrob/bonzai/pkg"

	"github.com/Chaitanyabsprip/dot/internal/core/embed"
	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

// TODO(me):
// - Create options to control what level of setup is done
// - The options will be one of slim/quik/full
// - Slim: install only the files that are needed and skip the rest
// - Quik: install only the configuration files and skip the
//   optional dependencies
// - Full: install everything, all config files and all dependencies

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Usage: `tmux setup`,
	Opts:  `slim|full`,
	Short: `Setup tmux copies configuration files to config directory`,
	Long:  ``,
	Call: func(x *bonzai.Cmd, _ ...string) error {
		return embed.SetupAll(embedFs, "tmux", oscfg.ConfigDir())
	},
}
