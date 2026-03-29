package vimium

import (
	"embed"
	"fmt"
	"path"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed vimium_c.json
var embedFs embed.FS

var SetupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Short: `setup vimium`,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		return e.SetupAll(embedFs, "vimium", oscfg.ConfigDir(), nil)
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit vimium configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(
			oscfg.ConfigDir(),
			"vimium_c.json",
		)
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run vimium setup")
		return nil
	},
}
