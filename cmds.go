package dot

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
	"github.com/rwxrob/bonzai/pkg/core/comp"

	"github.com/Chaitanyabsprip/dot/internal/alacritty"
	"github.com/Chaitanyabsprip/dot/internal/gh"
	"github.com/Chaitanyabsprip/dot/internal/git"
	"github.com/Chaitanyabsprip/dot/internal/tmux"
)

var Cmd = &bonzai.Cmd{
	Name:  `dot`,
	Alias: `d`,
	Usage: `dotfiles`,
	Short: `Manage dotfiles`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{tmux.Cmd, gh.Cmd, git.Cmd, alacritty.Cmd},
}
