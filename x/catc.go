package x

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/run"
)

var catcCmd = &bonzai.Cmd{
	Name:    `catc`,
	Vers:    `v1.0.0`,
	Short:   `find a script in path and cat the contents`,
	MinArgs: 1,
	Comp:    comp.Opts,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			return errors.New("no arguments provided")
		}
		catc(args...)
		return nil
	},
}

func catc(names ...string) {
	for _, name := range names {
		path, err := exec.LookPath(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding %s: %v\n", name, err)
			continue
		}
		err = run.SysExec("bat", path)
		if err != nil {
			// fallback to cat if bat is not available
			err = run.SysExec("cat", path)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"Error displaying %s: %v\n",
					name,
					err,
				)
			}
		}
	}
}
