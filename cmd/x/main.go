// Package main is the entry point for the x command line utilities.
// It initializes and executes the main x command tree.
package main

import (
	"github.com/Chaitanyabsprip/dotfiles/x"
	"github.com/Chaitanyabsprip/dotfiles/x/work"
)

func main() {
	if work.IsRefreshWorker() {
		work.RefreshCacheWorker()
		return
	}
	x.Cmd.Exec()
}
