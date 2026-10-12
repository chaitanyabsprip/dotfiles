// Package mise provides functionality for managing mise, the runtime
// version manager (node, python, ...) that replaces nvm/pyenv/rbenv.
package mise

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/edit"
	"github.com/rwxrob/bonzai/run"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

//go:embed mise
var embedFs embed.FS

var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Alias: `i`,
	Short: `install mise and the runtimes in its config`,
	Do:    func(_ *bonzai.Cmd, _ ...string) error { return installMise() },
}

var SetupCmd = &bonzai.Cmd{
	Name: `setup`,

	Short: `setup mise`,
	Comp:  comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return e.SetupAll(embedFs, "mise", oscfg.ConfigDir(), nil)
	},
}

var EditCmd = &bonzai.Cmd{
	Name:   `edit`,
	Short:  `edit mise configuration`,
	NoArgs: true,
	Do: func(x *bonzai.Cmd, _ ...string) error {
		filePath := path.Join(oscfg.ConfigDir(), "mise", "config.toml")
		if err := edit.Files(filePath); err != nil {
			return err
		}
		fmt.Println("rebuild binary")
		fmt.Println("re run mise setup")
		return nil
	},
}

// installMise installs mise (Homebrew on macOS, the official installer
// elsewhere, as ~/dotfiles setup does), then the runtimes pinned in
// ~/.config/mise/config.toml.
func installMise() error {
	if ok, _ := have.Executable(`mise`); !ok {
		var err error
		if ok, _ := have.Executable(`brew`); ok {
			err = run.Exec(`brew`, `install`, `mise`)
		} else {
			err = curlInstall()
		}
		if err != nil {
			return err
		}
	}
	bin, err := exec.LookPath(`mise`)
	if err != nil {
		// mise.run installs to ~/.local/bin, which may not be on PATH yet.
		bin = filepath.Join(oscfg.BinDir(), `mise`)
	}
	return run.Exec(bin, `install`)
}

func curlInstall() error {
	curl := exec.Command(`curl`, `-sS`, `https://mise.run`)
	sh := exec.Command(`sh`)
	sh.Stdin, _ = curl.StdoutPipe()
	sh.Stdout, sh.Stderr = os.Stdout, os.Stderr
	if err := curl.Start(); err != nil {
		return err
	}
	if err := sh.Run(); err != nil {
		return err
	}
	return curl.Wait()
}
