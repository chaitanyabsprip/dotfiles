// Package icon provides the `icon` command which displays the icon
// and name for the foreground process of a tmux pane.
package icon

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai"
)

var Cmd = &bonzai.Cmd{
	Name:  `icon`,
	Alias: `iname|i|name`,
	Short: `display icon and name for pane foreground process`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			fmt.Fprint(os.Stderr, "?")
			return nil
		}
		return Icon(args[0])
	},
}

func Icon(panePid string) error {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, "?")
		return nil
	}

	normalized, err := ResolveCommand(panePid)
	if err != nil {
		fmt.Print(cfg.Config.FallbackIcon)
		return nil
	}

	entry := Match(cfg, normalized, normalized)
	output := Format(entry, cfg.Config, normalized)
	fmt.Print(output)
	return nil
}
