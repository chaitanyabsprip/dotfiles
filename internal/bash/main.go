package bash

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/internal/ohmyposh"
	"github.com/Chaitanyabsprip/dotfiles/internal/shell"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

//go:embed bashrc
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bash`,
	Short: `bash is a utility to manage bash configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup bash`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		shell.Cmd.Run("setup")
		home := os.Getenv("HOME")
		if err := shell.Cmd.Run(`setup`); err != nil {
			return err
		}
		if err := ohmyposh.Cmd.Run(`setup`); err != nil {
			return err
		}
		if err := install.OhMyPosh(); err != nil {
			return err
		}
		return e.SetupAll(
			embedFs,
			"bash",
			home,
			map[string]string{
				"bashrc": filepath.Join(home, ".bashrc"),
			},
		)
	},
}
