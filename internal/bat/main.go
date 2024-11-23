package bat

import (
	"embed"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/x/install"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed bat
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bat`,
	Short: `bat is a utility to manage bat configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, install.BatCmd.WithName(`install`)},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup bat`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, `bat`, oscfg.ConfigDir(), nil)
	},
}
