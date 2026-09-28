package colors

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

var Color255Cmd = &bonzai.Cmd{
	Name:  `ctwo`,
	Short: `print 256 colors in terminal`,
	Long: `
Prints all 256 terminal colors: their numbers as foreground text on
one line, then again as background swatches on the next.`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{},
	Do: func(x *bonzai.Cmd, args ...string) error {
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
