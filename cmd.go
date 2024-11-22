package dot

import (
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dot/internal/alacritty"
	"github.com/Chaitanyabsprip/dot/internal/bash"
	"github.com/Chaitanyabsprip/dot/internal/bat"
	"github.com/Chaitanyabsprip/dot/internal/bin"
	"github.com/Chaitanyabsprip/dot/internal/brew"
	"github.com/Chaitanyabsprip/dot/internal/dirs"
	"github.com/Chaitanyabsprip/dot/internal/fish"
	"github.com/Chaitanyabsprip/dot/internal/gh"
	"github.com/Chaitanyabsprip/dot/internal/git"
	"github.com/Chaitanyabsprip/dot/internal/gitui"
	"github.com/Chaitanyabsprip/dot/internal/hypr"
	"github.com/Chaitanyabsprip/dot/internal/kitty"
	"github.com/Chaitanyabsprip/dot/internal/lsd"
	"github.com/Chaitanyabsprip/dot/internal/ohmyposh"
	"github.com/Chaitanyabsprip/dot/internal/shell"
	"github.com/Chaitanyabsprip/dot/internal/sqlfluff"
	"github.com/Chaitanyabsprip/dot/internal/starship"
	"github.com/Chaitanyabsprip/dot/internal/tmux"
	"github.com/Chaitanyabsprip/dot/internal/vimium"
	"github.com/Chaitanyabsprip/dot/internal/waybar"
	"github.com/Chaitanyabsprip/dot/internal/zsh"
	"github.com/Chaitanyabsprip/dot/x"
)

var Cmd = &bonzai.Cmd{
	Name:  `dot`,
	Alias: `d`,
	Short: `Manage dotfiles`,
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
