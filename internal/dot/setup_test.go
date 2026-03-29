package dot

import (
	"testing"

	"github.com/rwxrob/bonzai"
)

func TestSetupCmdDoNoArgsRunsAll(t *testing.T) {
	orig := SetupCmds
	defer func() { SetupCmds = orig }()

	called := map[string]bool{}
	SetupCmds = []*bonzai.Cmd{
		{Name: `bash`, Do: func(x *bonzai.Cmd, _ ...string) error { called[`bash`] = true; return nil }},
		{Name: `zsh`, Do: func(x *bonzai.Cmd, _ ...string) error { called[`zsh`] = true; return nil }},
	}

	if err := SetupCmd.Do(SetupCmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, name := range []string{`bash`, `zsh`} {
		if !called[name] {
			t.Errorf("tool %q was not set up by no-args Do", name)
		}
	}
}

func TestSetupCmdDoSpecificTool(t *testing.T) {
	orig := SetupCmds
	defer func() { SetupCmds = orig }()

	called := false
	SetupCmds = []*bonzai.Cmd{
		{Name: `bash`, Do: func(x *bonzai.Cmd, _ ...string) error { called = true; return nil }},
	}

	if err := SetupCmd.Do(SetupCmd, `bash`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("bash setup was not called for Do('bash')")
	}
}
