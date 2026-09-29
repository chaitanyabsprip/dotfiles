// Package glow provides functionality for managing glow (markdown reader)
// configuration.
package glow

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

//go:embed glow
var embedFs embed.FS

var SetupCmd = &bonzai.Cmd{
	Name: `setup`,

	Short: `setup glow`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "glow", oscfg.ConfigDir(), nil)
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit glow configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(oscfg.ConfigDir(), "glow", "glow.yml")
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run glow setup")
		return nil
	},
}
