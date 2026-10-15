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

func StateDir() string {
	if len(env.XdgStateHome) > 0 {
		return env.XdgStateHome
	}
	return filepath.Join(env.Home, ".local", "state")
}

// BinDir is ~/.local/bin, the untracked dir installers drop binaries in.
func BinDir() string {
	return filepath.Join(env.Home, ".local", "bin")
}

// DotfilesDir is ~/dotfiles, where the shell, bin and oh-my-posh configs
// live: the embedded configs reference them there by literal path
// ("$HOME/dotfiles/shell/env.sh"), so it ignores $DOTFILES. On a machine
// with a ~/dotfiles git checkout it is that checkout.
func DotfilesDir() string {
	return filepath.Join(env.Home, "dotfiles")
}

func CacheDir() string {
	if cacheDir := env.XdfCacheHome; len(cacheDir) > 0 {
		return cacheDir
	}
	return filepath.Join(env.Home, ".cache")
}
