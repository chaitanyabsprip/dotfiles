package bin

import (
	"embed"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

//go:embed bin
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bin`,
	Short: `manage scripts`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup bin directory`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		binDir := filepath.Join(env.HOME, ".local")
		return e.SetupAll(embedFs, "bin", binDir, nil)
	},
}
