package embed

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Chaitanyabsprip/dot/internal/core/oscfg"
)

func SetupAll(
	embedFs embed.FS,
	name, configDir string,
	altPaths map[string]string,
) error {
	toolDir := filepath.Join(configDir, name)
	if _, err := os.Stat(toolDir); err == nil {
		err := os.RemoveAll(toolDir)
		if err != nil {
			return err
		}
		return CopyFiles(embedFs, name, configDir, altPaths)
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
	return CopyFiles(embedFs, name, configDir, altPaths)
}

func CopyFiles(
	embedFs embed.FS,
	name, configDir string,
	altPaths map[string]string,
) error {
	if len(configDir) == 0 {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return fs.WalkDir(
		embedFs,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			targetPath := filepath.Join(configDir, path)
			if altPath, ok := altPaths[path]; ok {
				targetPath = filepath.Join(
					altPath,
					filepath.Base(path),
				)
			}
			fmt.Printf(
				"path: %s, targetPath: %s\n",
				path,
				targetPath,
			)
			if d.IsDir() {
				return os.MkdirAll(targetPath, 0o755)
			}
			content, err := fs.ReadFile(embedFs, path)
			if err != nil {
				return err
			}
			os.WriteFile(
				targetPath,
				content,
				getFileMode(path),
			)
			return nil
		},
	)
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
