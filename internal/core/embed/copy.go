// Package embed provides utilities for working with embedded files and directories.
// It facilitates copying embedded filesystem contents to the host system,
// handling operations such as setup, configuration deployment, and file management
// with appropriate permissions and path handling.
package embed

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

const Skip = ""

func SetupAll(
	embedFs embed.FS,
	name, configDir string,
	overrides map[string]string,
) error {
	toolDir := filepath.Join(configDir, name)
	if _, err := os.Stat(toolDir); err == nil {
		err := os.RemoveAll(toolDir)
		if err != nil {
			return err
		}
		return CopyAllFiles(embedFs, name, configDir, overrides)
	}
	if f, err := os.Stat(toolDir); err == nil {
		if !f.IsDir() {
			err := os.Remove(toolDir)
			if err != nil {
				return err
			}
		}
		err := os.Rename(toolDir, oscfg.BackupDir(configDir))
		if err != nil {
			return err
		}
	}
	return CopyAllFiles(embedFs, name, configDir, overrides)
}

func CopyFilesRegx(
	embedFs embed.FS,
	name, configDir, pattern string,
	overrides map[string]string,
) error {
	if len(configDir) == 0 {
		configDir = filepath.Join(env.HOME, ".config")
	}
	regx, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	return fs.WalkDir(
		embedFs,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if regx.Match([]byte(path)) {
				targetPath := filepath.Join(
					configDir,
					path,
				)
				if altPath, ok := overrides[path]; ok {
					if altPath == Skip {
						return nil
					}
					targetPath = altPath
				}
				copy(embedFs, d, path, targetPath)
			}
			return nil
		},
	)
}

func CopyAllFiles(
	embedFs embed.FS,
	name, configDir string,
	overrides map[string]string,
) error {
	if len(configDir) == 0 {
		configDir = filepath.Join(env.HOME, ".config")
	}
	return fs.WalkDir(
		embedFs,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			targetPath := filepath.Join(configDir, path)
			if altPath, ok := overrides[path]; ok {
				if altPath == Skip {
					return nil
				}
				targetPath = altPath
			}
			fmt.Printf(
				"path: %s, targetPath: %s\n",
				path,
				targetPath,
			)
			err = copy(embedFs, d, path, targetPath)
			if err != nil {
				return err
			}
			return nil
		},
	)
}

func copy(embedFs embed.FS, d fs.DirEntry, path, dest string) error {
	if d.IsDir() {
		return os.MkdirAll(dest, 0o755)
	}
	content, err := fs.ReadFile(embedFs, path)
	if err != nil {
		return err
	}
	os.WriteFile(dest, content, getFileMode(path))
	return nil
}

func getFileMode(path string) fs.FileMode {
	var mode fs.FileMode
	if strings.Contains(path, "/bin/") ||
		strings.HasPrefix(path, "bin/") {
		mode = 0o755
	} else {
		mode = 0o644
	}
	return mode
}
