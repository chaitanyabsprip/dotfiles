package sqlfluff

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

//go:embed all:sqlfluff
var embedFs embed.FS

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Do:    func(_ *bonzai.Cmd, _ ...string) error { return installSqlfluff() },
}

var SetupCmd = &bonzai.Cmd{
	Name: `setup`,

	Short: `setup sqlfluff`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "sqlfluff", oscfg.ConfigDir(), nil)
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit sqlfluff configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(
			oscfg.ConfigDir(),
			"sqlfluff",
			".sqlfluff",
		)
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run sqlfluff setup")
		return nil
	},
}
