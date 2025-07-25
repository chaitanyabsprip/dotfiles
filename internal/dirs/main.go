package dirs

import (
	"embed"
	"fmt"
	"path"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed dirs
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `dirs`,
	Short: `manage dirs configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd, editCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup dirs`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		// TODO(chaitanya): install xdg-user-dirs
		return e.SetupAll(embedFs, "dirs", oscfg.ConfigDir(),
			map[string]string{
				`dirs`: filepath.Join(
					oscfg.ConfigDir(),
					`user-dirs.dirs`,
				),
			},
		)
	},
}

var editCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit dirs configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(oscfg.ConfigDir(), `user-dirs.dirs`)
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run dirs setup")
		return nil
	},
}
