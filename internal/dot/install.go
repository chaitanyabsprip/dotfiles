package dot

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/bash"
	"github.com/Chaitanyabsprip/dotfiles/internal/bat"
	"github.com/Chaitanyabsprip/dotfiles/internal/ohmyposh"
	"github.com/Chaitanyabsprip/dotfiles/internal/tmux"
	"github.com/Chaitanyabsprip/dotfiles/internal/zsh"
)

// InstallCmds contains per-tool install commands composed from tools that have them.
var InstallCmds = []*bonzai.Cmd{
	bash.InstallCmd.WithName(`bash`),
	bat.InstallCmd.WithName(`bat`),
	ohmyposh.InstallCmd.WithName(`ohmyposh`),
	zsh.InstallCmd.WithName(`zsh`),
	tmux.InstallCmd.WithName(`tmux`),
}

var allCmd = &bonzai.Cmd{
	Name:  `all`,
	Short: `install all tools`,
	Do:    func(_ *bonzai.Cmd, _ ...string) error { return installAll() },
}

func installAll() error {
	for _, cmd := range InstallCmds {
		if err := cmd.Run(""); err != nil {
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
				return cmd.Run(args[1:]...)
			}
		}
		return fmt.Errorf(`unknown tool: %s`, args[0])
	},
}
