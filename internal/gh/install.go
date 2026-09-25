package gh

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installGh() error {
	if ok, _ := have.Executable(`gh`); ok {
		fmt.Println(`gh is already installed`)
		return nil
	}
	return install.Pkg(`gh`, map[string]string{`pacman`: `github-cli`})
}
