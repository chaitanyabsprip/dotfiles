package starship

import (
	"embed"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed starship.toml
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `starship`,
	Short: `manage starship configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup starship`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "starship", oscfg.ConfigDir(), nil)
	},
}
