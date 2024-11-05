package bin

import (
	"embed"
	"os"
	"path/filepath"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai"
)

//go:embed bin
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bin`,
	Usage: `bin <command>`,
	Short: `bin is a utility to manage github-cli configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
	Init: func(x *bonzai.Cmd, args ...string) error {
		for _, cmd := range x.Cmds {
			cmd.Caller = x
		}
		return nil
	},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Usage: `setup <opts>`,
	Opts:  `slim|quik|full`,
	Short: `Setup bin directory`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		binDir := filepath.Join(os.Getenv("HOME"), ".local")
		return e.SetupAll(embedFs, "bin", binDir, nil)
	},
}
