package zsh

import (
	"embed"
	"fmt"
	"path"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/pkg/env"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed all:zsh
var embedFs embed.FS

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Do:    func(_ *bonzai.Cmd, _ ...string) error { return installZsh() },
}

var SetupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Alias: `conf`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		zshenvPath := filepath.Join(env.Home, `.zshenv`)
		overrides := map[string]string{
			`zsh/.zshenv`: zshenvPath,
		}
		err := e.SetupAll(
			embedFs,
			`zsh`,
			oscfg.ConfigDir(),
			overrides,
		)
		return err
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit zsh configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		cfgDir := oscfg.ConfigDir()
		filePath := path.Join(cfgDir, `zsh`)
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run zsh setup")
		return nil
	},
}
