package lsd

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installLsd() error {
	if ok, _ := have.Executable(`lsd`); ok {
		fmt.Println(`lsd is already installed`)
		return nil
	}
	return install.Pkg(`lsd`, nil)
}
