package waybar

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installWaybar() error {
	if ok, _ := have.Executable(`waybar`); ok {
		fmt.Println(`waybar is already installed`)
		return nil
	}
	// waybar is a Wayland status bar; there is no macOS build.
	return install.Pkg(`waybar`, map[string]string{`brew`: ``})
}
