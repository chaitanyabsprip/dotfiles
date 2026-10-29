// Package ohmyposh provides functionality for managing Oh My Posh shell prompt configuration.
// It offers commands to setup and customize the Oh My Posh prompt engine, allowing
// for customized shell prompts with enhanced styling and information display.
package ohmyposh

import (
	"embed"
	"fmt"
	"path"

	"github.com/chaitanyabsprip/bonzai"
	"github.com/chaitanyabsprip/bonzai/comp"
	"github.com/chaitanyabsprip/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed oh-my-posh.rc.toml
var embedFs embed.FS

var SetupCmd = &bonzai.Cmd{
	Name: `setup`,

	Short: `setup ohmyposh`,
	Comp:  comp.Opts,
	// ~/dotfiles/oh-my-posh.rc.toml: zsh, bash and fish pass
	// --config "$DOTFILES/oh-my-posh.rc.toml".
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "ohmyposh", oscfg.DotfilesDir(), nil)
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit ohmyposh configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(oscfg.DotfilesDir(), "oh-my-posh.rc.toml")
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run ohmyposh setup")
		return nil
	},
}
