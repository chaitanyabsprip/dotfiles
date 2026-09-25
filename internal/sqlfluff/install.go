package sqlfluff

import (
	"fmt"

	"github.com/Chaitanyabsprip/dotfiles/x/have"
	"github.com/Chaitanyabsprip/dotfiles/x/install"
)

func installSqlfluff() error {
	if ok, _ := have.Executable(`sqlfluff`); ok {
		fmt.Println(`sqlfluff is already installed`)
		return nil
	}
	return install.Pkg(`sqlfluff`, map[string]string{`dnf`: ``})
}
