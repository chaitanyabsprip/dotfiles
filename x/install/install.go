// Package install provides utilities for installing and managing software packages
// and tools. It includes functions for downloading files, extracting archives,
// and running commands with root privileges. The package supports installation
// from GitHub releases and other sources.
package install

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/github"
	"github.com/rwxrob/bonzai/run"
	"github.com/rwxrob/bonzai/web"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

func WithRoot(args ...string) error {
	if os.Geteuid() != 0 {
		if _, err := exec.LookPath(`sudo`); err != nil {
			return fmt.Errorf(
				`user not root and sudo not found`,
			)
		}
		args = append([]string{`sudo`}, args...)
	}
	return run.Exec(args...)
}

func GhDownload(
	repo, version, assetName string,
) (downloadPath string, err error) {
	fmtString := `https://%s/%s/releases/download/%s/%s`
	downloadPath = filepath.Join(BinDir, assetName)
	downloadUrl := fmt.Sprintf(
		fmtString,
		github.Host,
		repo,
		version,
		assetName,
	)
	err = DownloadFile(downloadUrl, downloadPath)
	return downloadPath, err
}

func DownloadFile(url, dest string) (err error) {
	file, err := os.Create(dest)
	if err != nil {
		return
	}
	defer func() { err = errors.Join(err, file.Close()) }()
	req := web.Req{U: url, D: file}
	err = req.Submit()
	return
}

func ExtractTarGz(tarPath, dest string) (err error) {
	if futil.NotExists(dest) {
		if err := os.MkdirAll(dest, 0o755); err != nil {
			return err
		}
		defer func() { // if ends with error, clean up
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
			break // End of archive
		}
		if err != nil {
			return err
		}

		path := filepath.Join(
			dest,
			removeFirstPart(header.Name),
		)
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
		fmt.Printf(
			"tar contains unsupported header type: %v\n",
			header.Typeflag,
		)
	}
	return nil
}

func removeFirstPart(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) <= 1 {
		return "" // or path, depending on what behavior you want
	}
	return filepath.Join(parts[1:]...)
}

var BinDir = filepath.Join(env.Home, `.local`, `bin`)
