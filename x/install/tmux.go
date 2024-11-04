package install

import (
	"fmt"
	"os"

	"github.com/Chaitanyabsprip/dot/x/distro"
	"github.com/Chaitanyabsprip/dot/x/have"
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/run"
)

var tmuxCmd = &bonzai.Cmd{
	Name: `tmux`,
	Call: func(x *bonzai.Cmd, args ...string) error { return Tmux() },
}

func Tmux() error {
	if ok, _ := have.Executable("tmux"); ok {
		fmt.Println("tmux is already installed")
		return nil
	}
	switch distro.Name() {
	case "Arch Linux":
		return WithRoot("pacman", "-S", "tmux")
	case "Ubuntu", "Debian":
		return WithRoot("apt-get", "install", "tmux")
	case "Fedora":
		return run.SysExec("dnf", "install", "tmux", "-y")
	case "Darwin":
		return run.SysExec("brew", "install", "tmux")
	default:
		fmt.Fprintln(
			os.Stderr,
			"Unsupported operating system. Please install tmux manually.",
		)
	}
	return nil
}
