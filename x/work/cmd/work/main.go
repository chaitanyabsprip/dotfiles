package main

import (
	"github.com/Chaitanyabsprip/dotfiles/x/work"
)

func main() {
	if work.IsRefreshWorker() {
		work.RefreshCacheWorker()
		return
	}
	work.Cmd.Exec()
}
