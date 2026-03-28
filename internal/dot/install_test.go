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
		{Name: `mock1`, Do: func(x *bonzai.Cmd, _ ...string) error { called[`mock1`] = true; return nil }},
		{Name: `mock2`, Do: func(x *bonzai.Cmd, _ ...string) error { called[`mock2`] = true; return nil }},
	}

	if err := InstallCmd.Do(InstallCmd, `all`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, name := range []string{`mock1`, `mock2`} {
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
