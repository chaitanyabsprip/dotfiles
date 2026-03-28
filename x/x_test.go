package x_test

import (
	"testing"

	"github.com/Chaitanyabsprip/dotfiles/x"
)

func TestCmdContainsTmuxXCmd(t *testing.T) {
	found := false
	for _, c := range x.Cmd.Cmds {
		if c.Name == "tmux" {
			found = true
			break
		}
	}
	if !found {
		t.Error("x.Cmd does not contain tmux subcommand")
	}
}

func TestCmdContainsStandaloneUtilities(t *testing.T) {
	utilities := []string{"pem", "case", "catc", "color", "creash", "depends", "distro", "gpt", "have", "work"}
	names := make(map[string]bool)
	for _, c := range x.Cmd.Cmds {
		names[c.Name] = true
	}
	for _, u := range utilities {
		if !names[u] {
			t.Errorf("x.Cmd does not contain utility %q", u)
		}
	}
}

func TestCmdDoesNotContainInstall(t *testing.T) {
	for _, c := range x.Cmd.Cmds {
		if c.Name == "install" {
			t.Error("x.Cmd must not contain 'install' subcommand")
		}
	}
}
