package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Chaitanyabsprip/dotfiles/x/workdirs"
)

func main() {
	short := flag.Bool("s", false, "short")
	refresh := flag.Bool("r", false, "refresh cache")
	flag.Parse()
	dirs := workdirs.CachedAllDirs(*refresh)
	if *short {
		fmt.Println(strings.Join(workdirs.Shorten(dirs), "\n"))
		return
	}
	fmt.Println(strings.Join(dirs, "\n"))
}
