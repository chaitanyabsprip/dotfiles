package dot

import "github.com/rwxrob/bonzai"

var InitCmd = &bonzai.Cmd{
	Name:    `init`,
	Short:   `bootstrap machine`,
	NumArgs: 0,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if err := setupAll(); err != nil {
			return err
		}
		return installAll()
	},
}
