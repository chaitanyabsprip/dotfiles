package depends

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/rwxrob/bonzai/term"
)

func On(onError func(error), dependencies ...string) {
	for _, dep := range dependencies {
		if !isInPath(dep) {
			if onError == nil {
				onError = DefaultOnError
			}
			onError(fmt.Errorf("program depends on %s", dep))
		}
	}
}

func DefaultOnError(err error) {
	fmt.Printf("%s, please install and try again.\n", err)
	if !term.IsInteractive() {
		syscall.Kill(os.Getppid(), syscall.SIGTERM)
	}
}

func isInPath(dep string) bool {
	_, err := exec.LookPath(dep)
	return err == nil
}
