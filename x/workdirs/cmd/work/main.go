package main

import (
	"os"

	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dot/x/workdirs"
)

func main() {
	os.Setenv("SHELL", "bash")
	if len(os.Getenv(`DEBUG`)) > 0 {
		run.AllowPanic = true
	}
	workdirs.Cmd.Run()
}
