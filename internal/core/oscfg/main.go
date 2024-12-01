package oscfg

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func ConfigDir() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if len(configDir) == 0 {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return configDir
}

func BackupDir(dir string) string {
	newName := fmt.Sprintf("%s.%d.old", dir, time.Now().UnixMilli())
	return filepath.Join(dir, newName)
}

func BinDir() string {
	return filepath.Join(os.Getenv("HOME"), ".local", "bin")
}
