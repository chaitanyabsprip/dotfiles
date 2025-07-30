// Package bash provides commands to manage bash configuration.
package bash

import (
	"embed"
	"fmt"
	"path"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/internal/ohmyposh"
	"github.com/Chaitanyabsprip/dotfiles/internal/shell"
	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

//go:embed bashrc
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bash`,
	Short: `bash is a utility to manage bash configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, editCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup bash`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		shell.Cmd.Run(`setup`)
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
			`bash`,
			env.Home,
			map[string]string{
				`bashrc`: filepath.Join(env.Home, `.bashrc`),
			},
		)
	},
}

var editCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit bash configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(env.Home, `.bashrc`)
		err := edit.Files(filePath)
		if err != nil {
			return err
		}
		fmt.Println(`rebuild binary`)
		fmt.Println(`re run bash setup`)
		return nil
	},
}
