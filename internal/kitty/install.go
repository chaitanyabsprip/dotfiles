package kitty

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installKitty() error {
	if ok, _ := have.Executable(`kitty`); ok {
		fmt.Println(`kitty is already installed`)
		return nil
	}
	return install.Pkg(`kitty`, map[string]string{`brew`: `cask:kitty`})
}
