package main

import (
	"github.com/Chaitanyabsprip/dotfiles/x/workdirs"
)

func main() {
	if workdirs.IsRefreshWorker() {
		workdirs.RefreshCacheWorker()
		return
	}
	workdirs.Cmd.Exec()
}
