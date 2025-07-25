package ohmyposh

import (
	"embed"
	"fmt"
	"path"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed oh-my-posh.rc.toml
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `ohmyposh`,
	Short: `manage ohmyposh configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, editCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup ohmyposh`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "ohmyposh", oscfg.ConfigDir(), nil)
	},
}

var editCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit ohmyposh configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(
			oscfg.ConfigDir(),
			"ohmyposh",
			"oh-my-posh.rc.toml",
		)
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run ohmyposh setup")
		return nil
	},
}
