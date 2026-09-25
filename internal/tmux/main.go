package tmux

import (
	"embed"
	"fmt"
	"path"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/edit"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/x/depends"
)

// Deps lists tmux's real runtime dependencies, per ADR-0011 — fzf backs
// the sessionizer, session manager, and notes picker (pkg/fzf).
var Deps = []depends.Dep{
	{Name: `tmux`, Note: `binary`},
	{Name: `fzf`, Note: `sessionizer`},
}

// TODO(me):
// - dependencies
// - live reload

//go:embed all:tmux
var embedFs embed.FS

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Do:    func(_ *bonzai.Cmd, _ ...string) error { return installTmux() },
}

var EditCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit tmux configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(oscfg.ConfigDir(), "tmux", "tmux.conf")
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run tmux setup")
		return nil
	},
}

var SetupCmd = &bonzai.Cmd{
	Name:  `init`,
	Short: `setup tmux (full setup)`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "tmux", oscfg.ConfigDir(), nil)
	},
}
