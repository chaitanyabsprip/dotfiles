package bash

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"
)

//go:embed bashrc
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bash`,
	Usage: `bash <command>`,
	Short: `bash is a utility to manage bash configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Usage: `setup <opts>`,
	Opts:  `slim|quik|full`,
	Short: `Setup bash`,
	Comp:  comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(
			embedFs,
			"bash",
			filepath.Join(os.Getenv("HOME"), `.bashrc`),
			nil,
		)
	},
}
