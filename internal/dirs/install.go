package dirs

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installDirs() error {
	if ok, _ := have.Executable(`xdg-user-dirs-update`); ok {
		fmt.Println(`xdg-user-dirs is already installed`)
		return nil
	}
	// xdg-user-dirs is part of the freedesktop XDG spec; no macOS build.
	return install.Pkg(`xdg-user-dirs`, map[string]string{`brew`: ``})
}
