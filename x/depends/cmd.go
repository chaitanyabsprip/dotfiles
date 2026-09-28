// Package depends provides utilities for checking and managing program dependencies.
// It allows commands to verify if required executables exist in the system PATH and
// provides appropriate error handling when dependencies are missing.
package depends

import "github.com/rwxrob/bonzai"

var Cmd = &bonzai.Cmd{
	Name:  `depends`,
	Short: `check dependencies are installed`,
	Long: `
Checks that every named program is installed. If any is missing, it
prints an error naming it and, unless run interactively, sends SIGTERM
to the parent process (not just this command) — so a script can
'depends foo bar' as a guard clause.`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		On(nil, args...)
		return nil
	},
}
