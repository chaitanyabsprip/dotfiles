package x_test

import (
	"testing"

	"github.com/rwxrob/bonzai"

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
	utilities := []string{"pem", "case", "catc", "color", "creash", "depends", "distro", "gpt", "have", "url", "work"}
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

// TestCmdTreeValid catches bad Short/Usage/Name values anywhere in the
// tree before they surface as a runtime error.
func TestCmdTreeValid(t *testing.T) {
	var walk func(c *bonzai.Cmd)
	walk = func(c *bonzai.Cmd) {
		if err := c.Validate(); err != nil {
			t.Errorf("%s: %v", c.Name, err)
		}
		for _, sub := range c.Cmds {
			walk(sub)
		}
	}
	walk(x.Cmd)
}

func TestCmdDoesNotContainInstall(t *testing.T) {
	for _, c := range x.Cmd.Cmds {
		if c.Name == "install" {
			t.Error("x.Cmd must not contain 'install' subcommand")
		}
	}
}
