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
	dirs := make([]string, 0)
	dirs = append(dirs, workdirs.Workdirs()...)
	dirs = append(dirs, workdirs.Worktrees()...)
	if *short {
		fmt.Println(strings.Join(workdirs.Shorten(dirs), "\n"))
		return
	}
	fmt.Println(strings.Join(dirs, "\n"))
}
