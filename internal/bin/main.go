package bin

import (
	"embed"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
)

//go:embed bin
var embedFs embed.FS

var Cmd = &bonzai.Cmd{
	Name:  `bin`,
	Short: `bin is a utility to manage github-cli configuration`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{setupCmd},
}

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Opts:  `slim|quik|full`,
	Short: `setup bin directory`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		binDir := filepath.Join(os.Getenv("HOME"), ".local")
		return e.SetupAll(embedFs, "bin", binDir, nil)
	},
}
