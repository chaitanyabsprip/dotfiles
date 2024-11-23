package dirs

import (
	"embed"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dot/internal/core/embed"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

//go:embed dirs
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `dirs`,
	Short: `dirs is a utility to manage dirs configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
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
