// Package bat provides commands to manage the bat configuration.
package bat

import (
	"embed"
	"fmt"
	"path"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"
	"github.com/rwxrob/bonzai/run"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/pkg/with"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
)

//go:embed bat
var embedFs embed.FS

var InstallCmd = &bonzai.Cmd{
	Name: `install`,
	Cmds: []*bonzai.Cmd{batGhInstallCmd, batPkgInstallCmd},
	Do:   func(x *bonzai.Cmd, args ...string) error { return installBat() },
}

var batPkgInstallCmd = &bonzai.Cmd{
	Name: `pkg`,
	Do:   func(x *bonzai.Cmd, args ...string) error { return batPkgInstall() },
}

var batGhInstallCmd = &bonzai.Cmd{
	Name: `gh`,
	Do:   func(x *bonzai.Cmd, args ...string) error { return batGhInstall() },
}

var SetupCmd = &bonzai.Cmd{
	Name: `setup`,

	Short: `setup bat`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) (err error) {
		err = e.SetupAll(embedFs, `bat`, oscfg.ConfigDir(), nil)
		if err != nil || e.Check != e.CheckOff {
			// `dot status`/`dot diff` run every tool's SetupCmd in a
			// read-only check mode; skip the real `bat cache --build` side
			// effect below when that's what's happening.
			return err
		}
		reset, err := with.Path(oscfg.BinDir())
		if err != nil {
			return err
		}
		defer func() { err = reset() }()
		return run.Exec(`bat`, `cache`, `--build`)
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   "edit",
	Short:  `edit bat configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		cfgDir := oscfg.ConfigDir()
		filePath := path.Join(cfgDir, "bat", "config")
		err := edit.Files(filePath)
		if err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run bat setup")
		return nil
	},
}
