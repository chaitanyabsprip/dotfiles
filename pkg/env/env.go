// Package env provides access to important environment variables used throughout
// the dotfiles management system. It initializes default values for critical
// paths when environment variables are not set.
package env

import (
	"os"
	"path/filepath"
)

var (
	Dotfiles        = os.Getenv(`DOTFILES`)
	Downloads       = os.Getenv(`DOWNLOADS`)
	Editor          = os.Getenv(`EDITOR`)
	Home            = os.Getenv(`HOME`)
	Notespath       = os.Getenv(`NOTESPATH`)
	Path            = os.Getenv(`PATH`)
	Projects        = os.Getenv(`PROJECTS`)
	Pictures        = os.Getenv(`PICTURES`)
	Programs        = os.Getenv(`PROGRAMS`)
	Scripts         = os.Getenv(`SCRIPTS`)
	Tmux            = os.Getenv(`TMUX`)
	Verbose         = os.Getenv(`VERBOSE`)
	Visual          = os.Getenv(`VISUAL`)
	XdgConfigHome   = os.Getenv(`XDG_CONFIG_HOME`)
	XdfCacheHome    = os.Getenv(`XDG_CACHE_HOME`)
)

func init() {
	if len(Projects) == 0 {
		Projects = filepath.Join(Home, "projects")
	}
	if len(Programs) == 0 {
		Programs = filepath.Join(Home, "programs")
	}
	if len(Pictures) == 0 {
		Pictures = filepath.Join(Home, "pictures")
	}
}
