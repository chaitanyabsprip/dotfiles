package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Chaitanyabsprip/dot/x/workdirs"
)

func main() {
	short := flag.Bool("s", false, "short")
	flag.Parse()
	if *short {
		fmt.Println(
			strings.Join(workdirs.Shorten(workdirs.Worktrees()), "\n"),
		)
		return
	}
	fmt.Println(strings.Join(workdirs.Worktrees(), "\n"))
}
