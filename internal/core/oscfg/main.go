// Package oscfg provides functions to manage OS configuration
// directories and files.
package oscfg

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

func ConfigDir() string {
	configDir := env.XdgConfigHome
	if len(configDir) == 0 {
		configDir = filepath.Join(env.Home, ".config")
	}
	return configDir
}

func BackupDir(dir string) string {
	newName := fmt.Sprintf("%s.%d.old", dir, time.Now().UnixMilli())
	return filepath.Join(dir, newName)
}

func BinDir() string {
	return filepath.Join(env.Home, ".local", "bin")
}

func CacheDir() string {
	if cacheDir := env.XdfCacheHome; len(cacheDir) > 0 {
		return cacheDir
	}
	return filepath.Join(env.Home, ".cache")
}
