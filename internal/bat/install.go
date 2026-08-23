package bat

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/github"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installBat() error {
	if ok, _ := have.Executable(`bat`); ok {
		fmt.Println(`bat is already installed`)
		return nil
	}
	if err := batPkgInstall(); err == nil {
		return nil
	}
	if err := batGhInstall(); err != nil {
		return err
	}
	return fmt.Errorf(`unable to install bat`)
}

func batPkgInstall() error {
	switch distro.Name() {
	case `Arch Linux`:
		return install.WithRoot(`pacman`, `-S`, `--noconfirm`, `bat`)
	case `Ubuntu`, `Debian GNU/Linux`:
		err := install.WithRoot(`apt-get`, `install`, `-y`, `bat`)
		if err != nil {
			return err
		}
		binDir := oscfg.BinDir()
		if !futil.Exists(binDir) {
			if err := futil.CreateDir(binDir); err != nil {
				return err
			}
		}
		batPath := filepath.Join(binDir, `bat`)
		batcatPath, err := exec.LookPath(`batcat`)
		if err != nil {
			return err
		}
		if err := os.Symlink(batcatPath, batPath); err != nil {
			return err
		}
	case `Fedora`:
		return run.Exec(`dnf`, `install`, `bat`, `-y`)
	case `Darwin`:
		return run.Exec(`brew`, `install`, `tmux`)
	default:
		return fmt.Errorf(
			`unsupported or unconfigured operating system`,
		)
	}
	return nil
}

func batGhInstall() error {
	name, err := github.Latest(`sharkdp/bat`)
	if err != nil {
		return err
	}
	tarname := getBatTarname()
	binDir := oscfg.BinDir()
	if !futil.Exists(binDir) {
		if err := futil.CreateDir(binDir); err != nil {
			return err
		}
	}
	downloadPath, err := install.GhDownload(
		`sharkdp/bat`,
		name,
		tarname,
	)
	if err != nil {
		return err
	}
	extractPath := filepath.Join(binDir, `d-bat`)
	if err := install.ExtractTarGz(downloadPath, extractPath); err != nil {
		fmt.Println(err)
		return err
	}
	if err := os.Remove(downloadPath); err != nil {
		return err
	}
	srcPath := filepath.Join(extractPath, `bat`)
	destPath := filepath.Join(binDir, `bat`)
	if err := os.Rename(srcPath, destPath); err != nil {
		return err
	}
	if err := os.Chmod(destPath, 0o755); err != nil {
		return err
	}
	if err := os.RemoveAll(extractPath); err != nil {
		return err
	}
	return nil
}

func getBatTarname() string {
	switch fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH) {
	case `linux amd64`:
		return `bat-v0.24.0-x86_64-unknown-linux-musl.tar.gz`
	case `linux arm64`:
		return `bat-v0.24.0-aarch64-unknown-linux-gnu.tar.gz`
	case `darwin x86_64`:
		return `bat-v0.24.0-x86_64-apple-darwin.tar.gz`
	}
	return ``
}
