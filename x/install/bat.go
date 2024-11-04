package install

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/github"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dot/x/distro"
	"github.com/Chaitanyabsprip/dot/x/have"
)

var batCmd = &bonzai.Cmd{
	Name: `bat`,
	Call: func(x *bonzai.Cmd, args ...string) error { return Bat() },
}

func Bat() error {
	if ok, _ := have.Executable(`bat`); ok {
		fmt.Println(`bat is already installed`)
		// return nil
	}
	// if err := batPkgInstall(); err == nil {
	// 	return nil
	// }
	if err := batGhInstall(); err != nil {
		return err
	}
	return nil
}

func batPkgInstall() error {
	switch distro.Name() {
	case `Arch Linux`:
		WithRoot(`pacman`, `-S`, `bat`)
	case `Ubuntu`, `Debian`:
		return WithRoot(`apt-get`, `install`, `bat`)
	case `Fedora`:
		return run.SysExec(`dnf`, `install`, `bat`, `-y`)
	case `Darwin`:
		return run.SysExec(`brew`, `install`, `tmux`)
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
	downloadPath, err := GhDownload(`sharkdp/bat`, name, tarname)
	if err != nil {
		return err
	}
	extractPath := filepath.Join(BinDir, `d-bat`)
	if err := ExtractTarGz(downloadPath, extractPath); err != nil {
		fmt.Println(err)
		return err
	}
	if err := os.Remove(downloadPath); err != nil {
		return err
	}
	srcPath := filepath.Join(extractPath, `bat`)
	destPath := filepath.Join(BinDir, `bat`)
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
	tarname := ``
	switch fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH) {
	case `linux amd64`:
		tarname = `bat-v0.24.0-x86_64-unknown-linux-musl.tar.gz`
		// tarname =
		// `bat-v0.24.0-x86_64-unknown-linux-gnu.tar.gz`
	case `linux arm64`:
		tarname = `bat-v0.24.0-aarch64-unknown-linux-gnu.tar.gz`
	case `darwin x86_64`:
		tarname = `bat-v0.24.0-x86_64-apple-darwin.tar.gz`
	}
	return tarname
}
