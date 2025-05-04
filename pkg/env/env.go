package env

import "os"

var (
	DOTFILES        = os.Getenv(`DOTFILES`)
	DOWNLOADS       = os.Getenv(`DOWNLOADS`)
	EDITOR          = os.Getenv(`EDITOR`)
	HOME            = os.Getenv(`HOME`)
	NOTESPATH       = os.Getenv(`NOTESPATH`)
	PATH            = os.Getenv(`PATH`)
	PROJECTS        = os.Getenv(`PROJECTS`)
	SCRIPTS         = os.Getenv(`SCRIPTS`)
	TMUX            = os.Getenv(`TMUX`)
	VERBOSE         = os.Getenv(`VERBOSE`)
	VISUAL          = os.Getenv(`VISUAL`)
	XDG_CONFIG_HOME = os.Getenv(`XDG_CONFIG_HOME`)
)
