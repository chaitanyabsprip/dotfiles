package tmux

import (
	bonzai "github.com/rwxrob/bonzai/pkg"
)

var installCmd = &bonzai.Cmd{
	Name:  `install`,
	Usage: `tmux install`,
	Short: `Install tmux`,
}
