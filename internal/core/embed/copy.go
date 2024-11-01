package embed

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

func SetupAll(embedFs embed.FS, name, configDir string) error {
	toolDir := filepath.Join(configDir, name)
	setFilepath := filepath.Join(toolDir, ".set")
	if _, err := os.Stat(setFilepath); err == nil {
		err := os.RemoveAll(toolDir)
		if err != nil {
			return err
		}
		return CopyFiles(embedFs, name, configDir)
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
	return CopyFiles(embedFs, name, configDir)
}

func CopyFiles(embedFs embed.FS, name, configDir string) error {
	tmuxDir := filepath.Join(configDir, name)
	setFilepath := filepath.Join(tmuxDir, ".set")
	if len(configDir) == 0 {
		configDir = filepath.Join(
			os.Getenv("HOME"),
			".config",
		)
	}
	err := fs.WalkDir(
		embedFs,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			confPath := filepath.Join(configDir, path)
			if d.IsDir() {
				return os.MkdirAll(confPath, 0o755)
			}
			content, err := fs.ReadFile(embedFs, path)
			if err != nil {
				return err
			}
			os.WriteFile(
				confPath,
				content,
				getFileMode(path),
			)
			return nil
		},
	)
	if err != nil {
		return err
	}
	return os.WriteFile(setFilepath, []byte{}, 0o644)
}

func getFileMode(path string) fs.FileMode {
	var mode fs.FileMode
	if strings.Contains(path, "/bin/") {
		mode = 0o755
	} else {
		mode = 0o644
	}
	return mode
}
