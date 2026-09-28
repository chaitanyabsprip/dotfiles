package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Chaitanyabsprip/dotfiles/x/work"
)

func main() {
	if work.IsRefreshWorker() {
		work.RefreshCacheWorker()
		return
	}
	short := flag.Bool("s", false, "short")
	refresh := flag.Bool("r", false, "refresh cache")
	flag.Parse()
	dirs := work.CachedAllDirs(*refresh)
	if *short {
		fmt.Println(strings.Join(work.Shorten(dirs), "\n"))
		return
	}
	fmt.Println(strings.Join(dirs, "\n"))
}
