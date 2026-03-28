package dot

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/bash"
	"github.com/Chaitanyabsprip/dotfiles/internal/bat"
	"github.com/Chaitanyabsprip/dotfiles/internal/zsh"
)

// InstallCmds contains per-tool install commands composed from tools that have them.
var InstallCmds = []*bonzai.Cmd{
	{Name: `bash`, Do: func(x *bonzai.Cmd, args ...string) error { return bash.InstallCmd.Do(bash.InstallCmd, args...) }},
	{Name: `bat`, Do: func(x *bonzai.Cmd, args ...string) error { return bat.InstallCmd.Do(bat.InstallCmd, args...) }},
	{Name: `zsh`, Do: func(x *bonzai.Cmd, args ...string) error { return zsh.InstallCmd.Do(zsh.InstallCmd, args...) }},
}

var allCmd = &bonzai.Cmd{
	Name:  `all`,
	Short: `install all tools`,
	Do:    func(_ *bonzai.Cmd, _ ...string) error { return installAll() },
}

func installAll() error {
	for _, cmd := range InstallCmds {
		if err := cmd.Do(cmd); err != nil {
			return err
		}
	}
	return nil
}

// InstallCmd aggregates all tool install commands.
var InstallCmd = &bonzai.Cmd{
	Name:  `install`,
	Short: `install tools`,
	Comp:  comp.Cmds,
	Cmds:  append(InstallCmds, allCmd),
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			return fmt.Errorf(`specify a tool name or 'all'`)
		}
		if args[0] == `all` {
			return installAll()
		}
		for _, cmd := range InstallCmds {
			if cmd.Name == args[0] {
				return cmd.Do(cmd, args[1:]...)
			}
		}
		return fmt.Errorf(`unknown tool: %s`, args[0])
	},
}
