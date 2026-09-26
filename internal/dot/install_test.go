package dot

import (
	"testing"

	"github.com/rwxrob/bonzai"
)

func TestInstallCmdDoNoArgsReturnsError(t *testing.T) {
	if err := InstallCmd.Do(InstallCmd); err == nil {
		t.Error("InstallCmd.Do() with no args should return error, got nil")
	}
}

func TestInstallCmdDoAll(t *testing.T) {
	orig := InstallCmds
	defer func() { InstallCmds = orig }()

	called := map[string]bool{}
	InstallCmds = []*bonzai.Cmd{
		{Name: `bash`, Do: func(x *bonzai.Cmd, _ ...string) error { called[`bash`] = true; return nil }},
		{Name: `zsh`, Do: func(x *bonzai.Cmd, _ ...string) error { called[`zsh`] = true; return nil }},
	}

	if err := InstallCmd.Do(InstallCmd, `all`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, name := range []string{`bash`, `zsh`} {
		if !called[name] {
			t.Errorf("tool %q was not installed by 'all'", name)
		}
	}
}

func TestInstallCmdDoSpecificTool(t *testing.T) {
	orig := InstallCmds
	defer func() { InstallCmds = orig }()

	called := false
	InstallCmds = []*bonzai.Cmd{
		{Name: `bash`, Do: func(x *bonzai.Cmd, _ ...string) error { called = true; return nil }},
	}

	if err := InstallCmd.Do(InstallCmd, `bash`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("bash install was not called for Do('bash')")
	}
}

func TestInstallCmdsIncludeCLITools(t *testing.T) {
	names := map[string]bool{}
	for _, c := range InstallCmds {
		names[c.Name] = true
	}
	for _, want := range []string{`eza`, `fd`, `jq`, `neovim`, `rg`, `yq`} {
		if !names[want] {
			t.Errorf("dot install is missing %q", want)
		}
	}
}
