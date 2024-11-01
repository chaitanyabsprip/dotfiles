package dot

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/comp"

	"github.com/Chaitanyabsprip/dot/internal/alacritty"
	"github.com/Chaitanyabsprip/dot/internal/gh"
	"github.com/Chaitanyabsprip/dot/internal/git"
	"github.com/Chaitanyabsprip/dot/internal/gitui"
	"github.com/Chaitanyabsprip/dot/internal/kitty"
	"github.com/Chaitanyabsprip/dot/internal/tmux"
)

var Cmd = &bonzai.Cmd{
	Name:  `dot`,
	Alias: `d`,
	Usage: `dotfiles`,
	Short: `Manage dotfiles`,
	Comp:  comp.Cmds,
	Cmds: []*bonzai.Cmd{
		alacritty.Cmd,
		gh.Cmd,
		git.Cmd,
		gitui.Cmd,
		kitty.Cmd,
		tmux.Cmd,
	},
}
