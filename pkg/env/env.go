package env

import (
	"os"
	"path/filepath"
)

var (
	DOTFILES        = os.Getenv(`DOTFILES`)
	DOWNLOADS       = os.Getenv(`DOWNLOADS`)
	EDITOR          = os.Getenv(`EDITOR`)
	HOME            = os.Getenv(`HOME`)
	NOTESPATH       = os.Getenv(`NOTESPATH`)
	PATH            = os.Getenv(`PATH`)
	PROJECTS        = os.Getenv(`PROJECTS`)
	PICTURES        = os.Getenv(`PICTURES`)
	PROGRAMS        = os.Getenv(`PROGRAMS`)
	SCRIPTS         = os.Getenv(`SCRIPTS`)
	TMUX            = os.Getenv(`TMUX`)
	VERBOSE         = os.Getenv(`VERBOSE`)
	VISUAL          = os.Getenv(`VISUAL`)
	XDG_CONFIG_HOME = os.Getenv(`XDG_CONFIG_HOME`)
)

func init() {
	if len(PROJECTS) == 0 {
		PROJECTS = filepath.Join(HOME, "projects")
	}
	if len(PROGRAMS) == 0 {
		PROGRAMS = filepath.Join(HOME, "programs")
	}
	if len(PICTURES) == 0 {
		PICTURES = filepath.Join(HOME, "pictures")
	}
}
