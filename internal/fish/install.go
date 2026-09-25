package fish

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installFish() error {
	if ok, _ := have.Executable(`fish`); ok {
		fmt.Println(`fish is already installed`)
		return nil
	}
	return install.Pkg(`fish`, nil)
}
