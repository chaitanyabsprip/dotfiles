package ohmyposh

import (
	"embed"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed oh-my-posh.rc.toml
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `ohmyposh`,
	Usage: `ohmyposh <command>`,
	Short: `ohmyposh is a utility to manage ohmyposh configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Usage: `setup <opts>`,
	Opts:  `slim|quik|full`,
	Short: `Setup ohmyposh`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "ohmyposh", oscfg.ConfigDir(), nil)
	},
}
