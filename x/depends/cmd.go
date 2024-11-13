package depends

import "github.com/rwxrob/bonzai"

var Cmd = &bonzai.Cmd{
	Name: `depends`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		On(nil, args...)
		return nil
	},
}
