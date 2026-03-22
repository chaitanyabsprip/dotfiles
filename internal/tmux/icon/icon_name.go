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
		// Support both single argument (pane_pid only) and two arguments
		// (pane_pid and pane_current_command for fallback)
		var panePid, paneCommand string
		if len(args) >= 1 {
			panePid = args[0]
		}
		if len(args) >= 2 {
			paneCommand = args[1]
		}
		return Icon(panePid, paneCommand)
	},
}

func Icon(panePid string, paneCommand string) error {
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, "?")
		return nil
	}

	normalized, err := ResolveCommand(panePid)
	if err != nil {
		// No child process (idle pane) - use pane_current_command as fallback
		if paneCommand != "" {
			entry := Match(cfg, paneCommand, paneCommand)
			output := Format(entry, cfg.Config, paneCommand)
			fmt.Print(output)
		} else {
			fmt.Print(cfg.Config.FallbackIcon)
		}
		return nil
	}

	entry := Match(cfg, normalized, normalized)
	output := Format(entry, cfg.Config, normalized)
	fmt.Print(output)
	return nil
}

// IconSimple is kept for backward compatibility.
// For idle panes, it shows only the fallback icon.
func IconSimple(panePid string) error {
	return Icon(panePid, "")
}
