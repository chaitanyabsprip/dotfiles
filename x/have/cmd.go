package have

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai/comp"
	bonzai "github.com/rwxrob/bonzai/pkg"
)

var Cmd = &bonzai.Cmd{
	Name:  `have`,
	Usage: `have <command>`,
	Short: `have <command>`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
	Call: func(x *bonzai.Cmd, args ...string) error {
		verbose := os.Getenv("VERBOSE") != ""
		executables := os.Args[1:]
		for _, executable := range executables {
			ok, path := Have(executable)
			if !ok {
				if verbose {
					fmt.Fprintf(
						os.Stderr,
						"Executable not found: %s\n",
						executable,
					)
				}
				return fmt.Errorf(
					"executable not found: %s",
					executable,
				)
			} else if verbose {
				fmt.Printf("Found %s at %s\n", executable, path)
			}
		}
		return nil
	},
}

func Have(executable string) (bool, string) {
	out, err := exec.LookPath(executable)
	return err == nil, out
}
