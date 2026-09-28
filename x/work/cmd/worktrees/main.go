package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Chaitanyabsprip/dotfiles/x/work"
)

func main() {
	short := flag.Bool("s", false, "short")
	flag.Parse()
	if *short {
		fmt.Println(
			strings.Join(work.Shorten(work.Worktrees()), "\n"),
		)
		return
	}
	fmt.Println(strings.Join(work.Worktrees(), "\n"))
}
