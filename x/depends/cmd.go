// Package depends provides utilities for checking and managing program dependencies.
// It allows commands to verify if required executables exist in the system PATH and
// provides appropriate error handling when dependencies are missing.
package depends

import "github.com/rwxrob/bonzai"

var Cmd = &bonzai.Cmd{
	Name: `depends`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		On(nil, args...)
		return nil
	},
}
