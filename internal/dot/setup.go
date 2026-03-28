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

// SetupCmds contains per-tool setup commands composed from all 21 tools.
var SetupCmds = []*bonzai.Cmd{
	{Name: `alacritty`, Do: func(x *bonzai.Cmd, args ...string) error {
		return alacritty.SetupCmd.Do(alacritty.SetupCmd, args...)
	}},
	{Name: `bash`, Do: func(x *bonzai.Cmd, args ...string) error {
		return bash.SetupCmd.Do(bash.SetupCmd, args...)
	}},
	{Name: `bat`, Do: func(x *bonzai.Cmd, args ...string) error {
		return bat.SetupCmd.Do(bat.SetupCmd, args...)
	}},
	{Name: `bin`, Do: func(x *bonzai.Cmd, args ...string) error {
		return bin.SetupCmd.Do(bin.SetupCmd, args...)
	}},
	{Name: `brew`, Do: func(x *bonzai.Cmd, args ...string) error {
		return brew.SetupCmd.Do(brew.SetupCmd, args...)
	}},
	{Name: `dirs`, Do: func(x *bonzai.Cmd, args ...string) error {
		return dirs.SetupCmd.Do(dirs.SetupCmd, args...)
	}},
	{Name: `fish`, Do: func(x *bonzai.Cmd, args ...string) error {
		return fish.SetupCmd.Do(fish.SetupCmd, args...)
	}},
	{Name: `gh`, Do: func(x *bonzai.Cmd, args ...string) error {
		return gh.SetupCmd.Do(gh.SetupCmd, args...)
	}},
	{Name: `git`, Do: func(x *bonzai.Cmd, args ...string) error {
		return git.SetupCmd.Do(git.SetupCmd, args...)
	}},
	{Name: `gitui`, Do: func(x *bonzai.Cmd, args ...string) error {
		return gitui.SetupCmd.Do(gitui.SetupCmd, args...)
	}},
	{Name: `hypr`, Do: func(x *bonzai.Cmd, args ...string) error {
		return hypr.SetupCmd.Do(hypr.SetupCmd, args...)
	}},
	{Name: `kitty`, Do: func(x *bonzai.Cmd, args ...string) error {
		return kitty.SetupCmd.Do(kitty.SetupCmd, args...)
	}},
	{Name: `lsd`, Do: func(x *bonzai.Cmd, args ...string) error {
		return lsd.SetupCmd.Do(lsd.SetupCmd, args...)
	}},
	{Name: `ohmyposh`, Do: func(x *bonzai.Cmd, args ...string) error {
		return ohmyposh.SetupCmd.Do(ohmyposh.SetupCmd, args...)
	}},
	{Name: `shell`, Do: func(x *bonzai.Cmd, args ...string) error {
		return shell.SetupCmd.Do(shell.SetupCmd, args...)
	}},
	{Name: `sqlfluff`, Do: func(x *bonzai.Cmd, args ...string) error {
		return sqlfluff.SetupCmd.Do(sqlfluff.SetupCmd, args...)
	}},
	{Name: `starship`, Do: func(x *bonzai.Cmd, args ...string) error {
		return starship.SetupCmd.Do(starship.SetupCmd, args...)
	}},
	{Name: `tmux`, Do: func(x *bonzai.Cmd, args ...string) error {
		return tmux.SetupCmd.Do(tmux.SetupCmd, args...)
	}},
	{Name: `vimium`, Do: func(x *bonzai.Cmd, args ...string) error {
		return vimium.SetupCmd.Do(vimium.SetupCmd, args...)
	}},
	{Name: `waybar`, Do: func(x *bonzai.Cmd, args ...string) error {
		return waybar.SetupCmd.Do(waybar.SetupCmd, args...)
	}},
	{Name: `zsh`, Do: func(x *bonzai.Cmd, args ...string) error {
		return zsh.SetupCmd.Do(zsh.SetupCmd, args...)
	}},
}

func setupAll() error {
	for _, cmd := range SetupCmds {
		if err := cmd.Do(cmd); err != nil {
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
				return cmd.Do(cmd, args[1:]...)
			}
		}
		return fmt.Errorf(`unknown tool: %s`, args[0])
	},
}
