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
	"github.com/Chaitanyabsprip/dotfiles/internal/claude"
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

// SetupCmds contains per-tool setup commands composed from all 22 tools.
var SetupCmds = []*bonzai.Cmd{
	alacritty.SetupCmd.WithName(`alacritty`),
	bash.SetupCmd.WithName(`bash`),
	bat.SetupCmd.WithName(`bat`),
	bin.SetupCmd.WithName(`bin`),
	brew.SetupCmd.WithName(`brew`),
	claude.SetupCmd.WithName(`claude`),
	dirs.SetupCmd.WithName(`dirs`),
	fish.SetupCmd.WithName(`fish`),
	gh.SetupCmd.WithName(`gh`),
	git.SetupCmd.WithName(`git`),
	gitui.SetupCmd.WithName(`gitui`),
	hypr.SetupCmd.WithName(`hypr`),
	kitty.SetupCmd.WithName(`kitty`),
	lsd.SetupCmd.WithName(`lsd`),
	ohmyposh.SetupCmd.WithName(`ohmyposh`),
	shell.SetupCmd.WithName(`shell`),
	sqlfluff.SetupCmd.WithName(`sqlfluff`),
	starship.SetupCmd.WithName(`starship`),
	tmux.SetupCmd.WithName(`tmux`),
	vimium.SetupCmd.WithName(`vimium`),
	waybar.SetupCmd.WithName(`waybar`),
	zsh.SetupCmd.WithName(`zsh`),
}

func setupAll() error {
	for _, cmd := range SetupCmds {
		if err := cmd.Run(""); err != nil {
			return err
		}
	}
	return nil
}

// SetupCmd aggregates all tool setup commands.
var SetupCmd = &bonzai.Cmd{
	Name:  `setup`,
	Short: `setup tool configurations`,
	Comp:  comp.Cmds,
	Cmds:  SetupCmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			return setupAll()
		}
		for _, cmd := range SetupCmds {
			if cmd.Name == args[0] {
				return cmd.Run(args[1:]...)
			}
		}
		return fmt.Errorf(`unknown tool: %s`, args[0])
	},
}
