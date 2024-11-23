package colors

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

var TableCmd = &bonzai.Cmd{
	Name:  `table`,
	Short: `print color table`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
	Do: func(x *bonzai.Cmd, args ...string) error {
		Table()
		return nil
	},
}

// TODO(me): make this terminal size aware

func Table() {
	fmt.Println("Color escapes are \\e[${value};...;${value}m")
	fmt.Println("Values 30..37 are \033[33mforeground colors\033[m")
	fmt.Println("Values 40..47 are \033[43mbackground colors\033[m")
	fmt.Println("Value  1 gives a  \033[1mbold-faced look\033[m")
	fmt.Println()

	// Foreground colors
	for fgc := 30; fgc <= 37; fgc++ {
		// Background colors
		for bgc := 40; bgc <= 47; bgc++ {
			fgcText := ""
			bgcText := ""
			if fgc != 37 {
				fgcText = fmt.Sprintf("%d", fgc)
			}
			if bgc != 40 {
				bgcText = fmt.Sprintf("%d", bgc)
			}

			vals := ""
			if fgcText != "" && bgcText != "" {
				vals = fmt.Sprintf(
					"%s;%s",
					fgcText,
					bgcText,
				)
			} else if fgcText != "" {
				vals = fgcText
			} else if bgcText != "" {
				vals = bgcText
			}

			rSeq0 := fmt.Sprint("\\e[", vals, "m")
			fmt.Printf("  %-9s", rSeq0)
			seq0 := ""
			if vals != "" {
				seq0 = fmt.Sprint("\033[", vals, "m")
			}
			fmt.Printf(" %sTEXT\033[m", seq0)
			fmt.Printf(" \033[%s;1mBOLD\033[m", vals)
		}
		fmt.Println()
		fmt.Println()
	}
}
