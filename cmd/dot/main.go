// Package main is the entry point for the dot command line application.
// It initializes and executes the main dot command tree.
package main

import (
	dot "github.com/Chaitanyabsprip/dotfiles"
	"github.com/Chaitanyabsprip/dotfiles/x/work"
)

func main() {
	if work.IsRefreshWorker() {
		work.RefreshCacheWorker()
		return
	}
	dot.Cmd.Exec()
}
