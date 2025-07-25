// Package dot provides dotfile management functionality through a CLI application.
// It allows setting up, installing, and managing various configuration files for
// different tools and applications across the system.
package dot

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
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
	"github.com/Chaitanyabsprip/dotfiles/x"
)

var Cmd = &bonzai.Cmd{
	Name:  `dot`,
	Alias: `d`,
	Short: `manage dotfiles`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		alacritty.Cmd,
		bash.Cmd,
		bat.Cmd,
		bin.Cmd,
		brew.Cmd,
		dirs.Cmd,
		fish.Cmd,
		gh.Cmd,
		git.Cmd,
		gitui.Cmd,
		help.Cmd,
		hypr.Cmd,
		initCmd,
		kitty.Cmd,
		lsd.Cmd,
		ohmyposh.Cmd,
		shell.Cmd,
		sqlfluff.Cmd,
		starship.Cmd,
		tmux.Cmd,
		vimium.Cmd,
		waybar.Cmd,
		x.Cmd,
		zsh.Cmd,
	},
}

var initCmd = &bonzai.Cmd{
	Name:  `init`,
	Alias: `setup`,
	Short: `setup dotfiles`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		simpleCmds := []*bonzai.Cmd{
			bash.Cmd,
			bat.Cmd,
			bin.Cmd,
			brew.Cmd,
			dirs.Cmd,
			gh.Cmd,
			git.Cmd,
			gitui.Cmd,
			hypr.Cmd,
			kitty.Cmd,
			ohmyposh.Cmd,
			shell.Cmd,
			vimium.Cmd,
			waybar.Cmd,
		}
		for _, cmd := range simpleCmds {
			if err := cmd.Run(`setup`); err != nil {
				return err
			}
			if err := cmd.Run(`install`); err != nil {
				return err
			}
		}
		extCmds := []*bonzai.Cmd{tmux.Cmd, zsh.Cmd}
		for _, cmd := range extCmds {
			setupArgs := append([]string{`init`}, args...)
			if err := cmd.Run(setupArgs...); err != nil {
				return err
			}
		}
		return nil
	},
}
