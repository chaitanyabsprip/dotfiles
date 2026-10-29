// Package bash provides commands to manage bash configuration.
package bash

import (
	"embed"
	"fmt"
	"path"
	"path/filepath"

	"github.com/chaitanyabsprip/bonzai"
	"github.com/chaitanyabsprip/bonzai/comp"
	"github.com/chaitanyabsprip/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/internal/ohmyposh"
	"github.com/Chaitanyabsprip/dotfiles/internal/shell"
	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
	"github.com/Chaitanyabsprip/dotfiles/x/depends"
)

//go:embed bash
var embedFs embed.FS

var InstallCmd = &bonzai.Cmd{
	Name: `install`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if err := ohmyposh.SetupCmd.Run(``); err != nil {
			return err
		}
		return ohmyposh.InstallCmd.Run(``)
	},
}

// Deps mirrors InstallCmd above, which sets up and installs ohmyposh for
// the prompt — see ADR-0011.
var Deps = []depends.Dep{
	{Name: `ohmyposh`, Note: `prompt`, Deps: ohmyposh.Deps},
}

var SetupCmd = &bonzai.Cmd{
	Name: `setup`,

	Short: `setup bash`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if err := shell.SetupCmd.Run(``); err != nil {
			return err
		}
		return e.SetupAll(
			embedFs,
			`bash`,
			env.Home,
			// Machine-local lines belong in ~/.bash_profile.local, which
			// bash_profile sources and dot never touches.
			map[string]string{
				`bash`:              e.Skip, // don't create ~/bash
				`bash/bashrc`:       filepath.Join(env.Home, `.bashrc`),
				`bash/bash_profile`: filepath.Join(env.Home, `.bash_profile`),
			},
		)
	},
}

var EditCmd = &bonzai.Cmd{
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
