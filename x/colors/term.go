package colors

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

var TermCmd = &bonzai.Cmd{
	Name:  `term`,
	Short: `Print terminal colors`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
	Do: func(x *bonzai.Cmd, args ...string) error {
		Term()
		return nil
	},
}

func Term() {
	black := "\033[30m"
	red := "\033[31m"
	green := "\033[32m"
	yellow := "\033[33m"
	blue := "\033[34m"
	magenta := "\033[35m"
	cyan := "\033[36m"
	white := "\033[37m"
	reset := "\033[0m"

	fmt.Println("POSIX")
	fmt.Printf("%sblack=\"\\033[30m\"\n", black)
	fmt.Printf("%sred=\"\\033[31m\"\n", red)
	fmt.Printf("%sgreen=\"\\033[32m\"\n", green)
	fmt.Printf("%syellow=\"\\033[33m\"\n", yellow)
	fmt.Printf("%sblue=\"\\033[34m\"\n", blue)
	fmt.Printf("%smagenta=\"\\033[35m\"\n", magenta)
	fmt.Printf("%scyan=\"\\033[36m\"\n", cyan)
	fmt.Printf("%swhite=\"\\033[37m\"\n", white)
	fmt.Printf("reset=\"\\033[0m\"\n")

	fmt.Println()
	fmt.Println("\033[0mNC (No color)")
	fmt.Println("\033[1;37mWHITE\t\033[0;30mBLACK")
	fmt.Println("\033[0;34mBLUE\t\033[1;34mLIGHT_BLUE")
	fmt.Println("\033[0;32mGREEN\t\033[1;32mLIGHT_GREEN")
	fmt.Println("\033[0;36mCYAN\t\033[1;36mLIGHT_CYAN")
	fmt.Println("\033[0;31mRED\t\033[1;31mLIGHT_RED")
	fmt.Println("\033[0;35mPURPLE\t\033[1;35mLIGHT_PURPLE")
	fmt.Println("\033[0;33mYELLOW\t\033[1;33mLIGHT_YELLOW")
	fmt.Println("\033[1;30mGRAY\t\033[0;37mLIGHT_GRAY")
	fmt.Println(reset)
}
