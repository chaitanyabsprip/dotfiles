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
		hypr.Cmd,
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
		help.Cmd,
	},
}
