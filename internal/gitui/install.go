package gitui

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installGitui() error {
	if ok, _ := have.Executable(`gitui`); ok {
		fmt.Println(`gitui is already installed`)
		return nil
	}
	return install.Pkg(`gitui`, map[string]string{`apt`: ``})
}
