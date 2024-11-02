package colors

import (
	"fmt"

	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/comp"
)

var Color255Cmd = &bonzai.Cmd{
	Name:  `255`,
	Usage: `255`,
	Short: `Print 256 colors in terminal`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
	Call: func(x *bonzai.Cmd, args ...string) error {
		Color255()
		return nil
	},
}

func Color255() {
	for i := 0; i <= 255; i++ {
		fmt.Printf("\033[38;5;%dm%3d ", i, i)
	}
	fmt.Println("\033[0m")
	for i := 0; i <= 255; i++ {
		fmt.Printf("\033[48;5;%dm%3d ", i, i)
	}
	fmt.Println("\033[0m")
}
