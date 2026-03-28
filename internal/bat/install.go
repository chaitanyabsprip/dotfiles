package bat

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/github"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/web"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/x/distro"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
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
		return withRoot(`pacman`, `-S`, `bat`)
	case `Ubuntu`, `Debian GNU/Linux`:
		err := withRoot(`apt-get`, `install`, `-y`, `bat`)
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
		return fmt.Errorf(`unsupported or unconfigured operating system`)
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
	downloadPath, err := ghDownload(`sharkdp/bat`, name, tarname)
	if err != nil {
		return err
	}
	extractPath := filepath.Join(binDir, `d-bat`)
	if err := extractTarGz(downloadPath, extractPath); err != nil {
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

func ghDownload(repo, version, assetName string) (downloadPath string, err error) {
	downloadPath = filepath.Join(oscfg.BinDir(), assetName)
	downloadUrl := fmt.Sprintf(
		`https://%s/%s/releases/download/%s/%s`,
		github.Host, repo, version, assetName,
	)
	err = downloadFile(downloadUrl, downloadPath)
	return downloadPath, err
}

func downloadFile(url, dest string) (err error) {
	file, err := os.Create(dest)
	if err != nil {
		return
	}
	defer func() { err = errors.Join(err, file.Close()) }()
	req := web.Req{U: url, D: file}
	err = req.Submit()
	return
}

func extractTarGz(tarPath, dest string) (err error) {
	if futil.NotExists(dest) {
		if err := os.MkdirAll(dest, 0o755); err != nil {
			return err
		}
		defer func() {
			if err != nil {
				err = os.RemoveAll(dest)
			}
		}()
	}
	f, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, f.Close()) }()
	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		path := filepath.Join(dest, removeFirstPart(header.Name))
		if err := handleTarType(header, tarReader, dest, path); err != nil {
			return err
		}
	}
	return nil
}

func handleTarType(
	header *tar.Header,
	tarReader *tar.Reader,
	dest, path string,
) error {
	switch header.Typeflag {
	case tar.TypeDir:
		if err := os.MkdirAll(path, 0o755); err != nil {
			return err
		}
	case tar.TypeReg:
		outFile, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func() { err = errors.Join(err, outFile.Close()) }()
		if _, err := io.Copy(outFile, tarReader); err != nil {
			return err
		}
	case tar.TypeLink:
		linkPath := filepath.Join(dest, header.Linkname)
		if err := os.Link(linkPath, path); err != nil {
			return err
		}
	case tar.TypeSymlink:
		linkPath := filepath.Join(dest, header.Linkname)
		if err := os.Symlink(linkPath, path); err != nil {
			return err
		}
	default:
		fmt.Printf("tar contains unsupported header type: %v\n", header.Typeflag)
	}
	return nil
}

func removeFirstPart(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) <= 1 {
		return ""
	}
	return filepath.Join(parts[1:]...)
}

func withRoot(args ...string) error {
	if os.Geteuid() != 0 {
		if _, err := exec.LookPath(`sudo`); err != nil {
			return fmt.Errorf(`user not root and sudo not found`)
		}
		args = append([]string{`sudo`}, args...)
	}
	return run.Exec(args...)
}
