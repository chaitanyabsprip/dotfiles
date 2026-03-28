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
	{Name: `alacritty`, Do: func(x *bonzai.Cmd, args ...string) error { return alacritty.EditCmd.Do(alacritty.EditCmd, args...) }},
	{Name: `bash`, Do: func(x *bonzai.Cmd, args ...string) error { return bash.EditCmd.Do(bash.EditCmd, args...) }},
	{Name: `bat`, Do: func(x *bonzai.Cmd, args ...string) error { return bat.EditCmd.Do(bat.EditCmd, args...) }},
	{Name: `bin`, Do: func(x *bonzai.Cmd, args ...string) error { return bin.EditCmd.Do(bin.EditCmd, args...) }},
	{Name: `brew`, Do: func(x *bonzai.Cmd, args ...string) error { return brew.EditCmd.Do(brew.EditCmd, args...) }},
	{Name: `dirs`, Do: func(x *bonzai.Cmd, args ...string) error { return dirs.EditCmd.Do(dirs.EditCmd, args...) }},
	{Name: `fish`, Do: func(x *bonzai.Cmd, args ...string) error { return fish.EditCmd.Do(fish.EditCmd, args...) }},
	{Name: `gh`, Do: func(x *bonzai.Cmd, args ...string) error { return gh.EditCmd.Do(gh.EditCmd, args...) }},
	{Name: `git`, Do: func(x *bonzai.Cmd, args ...string) error { return git.EditCmd.Do(git.EditCmd, args...) }},
	{Name: `gitui`, Do: func(x *bonzai.Cmd, args ...string) error { return gitui.EditCmd.Do(gitui.EditCmd, args...) }},
	{Name: `hypr`, Do: func(x *bonzai.Cmd, args ...string) error { return hypr.EditCmd.Do(hypr.EditCmd, args...) }},
	{Name: `kitty`, Do: func(x *bonzai.Cmd, args ...string) error { return kitty.EditCmd.Do(kitty.EditCmd, args...) }},
	{Name: `lsd`, Do: func(x *bonzai.Cmd, args ...string) error { return lsd.EditCmd.Do(lsd.EditCmd, args...) }},
	{Name: `ohmyposh`, Do: func(x *bonzai.Cmd, args ...string) error { return ohmyposh.EditCmd.Do(ohmyposh.EditCmd, args...) }},
	{Name: `shell`, Do: func(x *bonzai.Cmd, args ...string) error { return shell.EditCmd.Do(shell.EditCmd, args...) }},
	{Name: `sqlfluff`, Do: func(x *bonzai.Cmd, args ...string) error { return sqlfluff.EditCmd.Do(sqlfluff.EditCmd, args...) }},
	{Name: `starship`, Do: func(x *bonzai.Cmd, args ...string) error { return starship.EditCmd.Do(starship.EditCmd, args...) }},
	{Name: `tmux`, Do: func(x *bonzai.Cmd, args ...string) error { return tmux.EditCmd.Do(tmux.EditCmd, args...) }},
	{Name: `vimium`, Do: func(x *bonzai.Cmd, args ...string) error { return vimium.EditCmd.Do(vimium.EditCmd, args...) }},
	{Name: `waybar`, Do: func(x *bonzai.Cmd, args ...string) error { return waybar.EditCmd.Do(waybar.EditCmd, args...) }},
	{Name: `zsh`, Do: func(x *bonzai.Cmd, args ...string) error { return zsh.EditCmd.Do(zsh.EditCmd, args...) }},
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
				return cmd.Do(cmd, args[1:]...)
			}
		}
		return fmt.Errorf(`unknown tool: %s`, args[0])
	},
}
