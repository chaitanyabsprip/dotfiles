package git

import (
	"embed"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed git
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `git`,
	Short: `git is a utility to manage git configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup git`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "git", oscfg.ConfigDir(), nil)
	},
}
