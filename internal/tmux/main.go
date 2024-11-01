package tmux

import (
	"embed"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/comp"
)

// TODO(me):
// - dependencies
// - live reload

//go:embed all:tmux
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `tmx`,
	Usage: `tmx <command>`,
	Short: `tmx is a utility to manage tmux configuration and related scripts`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, runCmd},
}
