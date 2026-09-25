package dot

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/alacritty"
	"github.com/Chaitanyabsprip/dotfiles/internal/bash"
	"github.com/Chaitanyabsprip/dotfiles/internal/bat"
	"github.com/Chaitanyabsprip/dotfiles/internal/dirs"
	"github.com/Chaitanyabsprip/dotfiles/internal/fish"
	"github.com/Chaitanyabsprip/dotfiles/internal/gh"
	"github.com/Chaitanyabsprip/dotfiles/internal/gitui"
	"github.com/Chaitanyabsprip/dotfiles/internal/kitty"
	"github.com/Chaitanyabsprip/dotfiles/internal/lsd"
	"github.com/Chaitanyabsprip/dotfiles/internal/ohmyposh"
	"github.com/Chaitanyabsprip/dotfiles/internal/sqlfluff"
	"github.com/Chaitanyabsprip/dotfiles/internal/starship"
	"github.com/Chaitanyabsprip/dotfiles/internal/tmux"
	"github.com/Chaitanyabsprip/dotfiles/internal/waybar"
	"github.com/Chaitanyabsprip/dotfiles/internal/zsh"
)

// InstallCmds contains per-tool install commands composed from tools that have them.
var InstallCmds = []*bonzai.Cmd{
	alacritty.InstallCmd.WithName(`alacritty`),
	bash.InstallCmd.WithName(`bash`),
	bat.InstallCmd.WithName(`bat`),
	dirs.InstallCmd.WithName(`dirs`),
	fish.InstallCmd.WithName(`fish`),
	gh.InstallCmd.WithName(`gh`),
	gitui.InstallCmd.WithName(`gitui`),
	kitty.InstallCmd.WithName(`kitty`),
	lsd.InstallCmd.WithName(`lsd`),
	ohmyposh.InstallCmd.WithName(`ohmyposh`),
	sqlfluff.InstallCmd.WithName(`sqlfluff`),
	starship.InstallCmd.WithName(`starship`),
	tmux.InstallCmd.WithName(`tmux`),
	waybar.InstallCmd.WithName(`waybar`),
	zsh.InstallCmd.WithName(`zsh`),
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
