package tmux

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	bonzai "github.com/rwxrob/bonzai/pkg"
)

// TODO(me):
// - Create options to control what level of setup is done
// - The options will be one of slim/quik/full
// - Slim: install only the files that are needed and skip the rest
// - Quik: install only the configuration files and skip the
//   optional dependencies
// - Full: install everything, all config files and all dependencies

var setupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Usage: `tmux setup`,
	Opts:  `slim|full`,
	Short: `Setup tmux copies configuration files to config directory`,
	Long:  ``,
	Call: func(x *bonzai.Cmd, _ ...string) error {
		configDir := getConfigDir()
		tmuxDir := filepath.Join(configDir, "tmux")
		setFilepath := filepath.Join(tmuxDir, ".set")
		if _, err := os.Stat(setFilepath); err == nil {
			err := os.RemoveAll(tmuxDir)
			if err != nil {
				return err
			}
			return copyFiles()
		}
		if f, err := os.Stat(tmuxDir); err == nil {
			if !f.IsDir() {
				err := os.Remove(tmuxDir)
				if err != nil {
					return err
				}
			}
			err := os.Rename(tmuxDir, getBackupDir())
			if err != nil {
				return err
			}
		}
		return copyFiles()
	},
}

func getBackupDir() string {
	configDir := getConfigDir()
	newName := fmt.Sprintf("tmux.%d.old", time.Now().UnixMilli())
	return filepath.Join(configDir, newName)
}

func copyFiles() error {
	configDir := getConfigDir()
	if len(configDir) == 0 {
		configDir = filepath.Join(
			os.Getenv("HOME"),
			".config",
		)
	}
	return fs.WalkDir(
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

func getConfigDir() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if len(configDir) == 0 {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return configDir
}
