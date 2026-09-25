package alacritty

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installAlacritty() error {
	if ok, _ := have.Executable(`alacritty`); ok {
		fmt.Println(`alacritty is already installed`)
		return nil
	}
	return install.Pkg(`alacritty`, map[string]string{`brew`: `cask:alacritty`})
}
