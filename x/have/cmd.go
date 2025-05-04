package have

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

var Cmd = &bonzai.Cmd{
	Name:  `have`,
	Short: `have <command>`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{},
	Do: func(x *bonzai.Cmd, args ...string) error {
		verbose := env.VERBOSE != ""
		names := os.Args[1:]
		for _, name := range names {
			ok, path := Executable(name)
			if !ok {
				if verbose {
					fmt.Fprintf(
						os.Stderr,
						"Executable not found: %s\n",
						name,
					)
				}
				return fmt.Errorf(
					"executable not found: %s",
					name,
				)
			} else if verbose {
				fmt.Printf("Found %s at %s\n", name, path)
			}
		}
		return nil
	},
}

func Executable(name string) (bool, string) {
	out, err := exec.LookPath(name)
	return err == nil, out
}
