package dot

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/alacritty"
	"github.com/Chaitanyabsprip/dotfiles/internal/bash"
	"github.com/Chaitanyabsprip/dotfiles/internal/bat"
	"github.com/Chaitanyabsprip/dotfiles/internal/bin"
	"github.com/Chaitanyabsprip/dotfiles/internal/brew"
	"github.com/Chaitanyabsprip/dotfiles/internal/dirs"
	"github.com/Chaitanyabsprip/dotfiles/internal/fish"
	"github.com/Chaitanyabsprip/dotfiles/internal/gh"
	"github.com/Chaitanyabsprip/dotfiles/internal/git"
	"github.com/Chaitanyabsprip/dotfiles/internal/gitui"
	"github.com/Chaitanyabsprip/dotfiles/internal/hypr"
	"github.com/Chaitanyabsprip/dotfiles/internal/kitty"
	"github.com/Chaitanyabsprip/dotfiles/internal/lsd"
	"github.com/Chaitanyabsprip/dotfiles/internal/ohmyposh"
	"github.com/Chaitanyabsprip/dotfiles/internal/shell"
	"github.com/Chaitanyabsprip/dotfiles/internal/sqlfluff"
	"github.com/Chaitanyabsprip/dotfiles/internal/starship"
	"github.com/Chaitanyabsprip/dotfiles/internal/tmux"
	"github.com/Chaitanyabsprip/dotfiles/internal/vimium"
	"github.com/Chaitanyabsprip/dotfiles/internal/waybar"
	"github.com/Chaitanyabsprip/dotfiles/internal/zsh"
)

// EditCmds contains per-tool edit commands composed from all tool packages.
var EditCmds = []*bonzai.Cmd{
	alacritty.EditCmd.WithName(`alacritty`),
	bash.EditCmd.WithName(`bash`),
	bat.EditCmd.WithName(`bat`),
	bin.EditCmd.WithName(`bin`),
	brew.EditCmd.WithName(`brew`),
	dirs.EditCmd.WithName(`dirs`),
	fish.EditCmd.WithName(`fish`),
	gh.EditCmd.WithName(`gh`),
	git.EditCmd.WithName(`git`),
	gitui.EditCmd.WithName(`gitui`),
	hypr.EditCmd.WithName(`hypr`),
	kitty.EditCmd.WithName(`kitty`),
	lsd.EditCmd.WithName(`lsd`),
	ohmyposh.EditCmd.WithName(`ohmyposh`),
	shell.EditCmd.WithName(`shell`),
	sqlfluff.EditCmd.WithName(`sqlfluff`),
	starship.EditCmd.WithName(`starship`),
	tmux.EditCmd.WithName(`tmux`),
	vimium.EditCmd.WithName(`vimium`),
	waybar.EditCmd.WithName(`waybar`),
	zsh.EditCmd.WithName(`zsh`),
}

// EditCmd aggregates all tool edit commands under a single verb command.
var EditCmd = &bonzai.Cmd{
	Name:  `edit`,
	Short: `edit tool configuration`,
	Comp:  comp.Cmds,
	Cmds:  EditCmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			return fmt.Errorf(`specify a tool to edit`)
		}
		for _, cmd := range EditCmds {
			if cmd.Name == args[0] {
				return cmd.Run(args[1:]...)
			}
		}
		return fmt.Errorf(`unknown tool: %s`, args[0])
	},
}
