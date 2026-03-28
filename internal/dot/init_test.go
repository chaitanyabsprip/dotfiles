package dot

import (
	"testing"

	"github.com/rwxrob/bonzai"
)

func TestInitCmdDoCallsSetupThenInstallAll(t *testing.T) {
	origSetup := SetupCmds
	origInstall := InstallCmds
	defer func() {
		SetupCmds = origSetup
		InstallCmds = origInstall
	}()

	var order []string

	SetupCmds = []*bonzai.Cmd{
		{Name: `mock`, Do: func(x *bonzai.Cmd, _ ...string) error {
			order = append(order, `setup`)
			return nil
		}},
	}
	InstallCmds = []*bonzai.Cmd{
		{Name: `mock`, Do: func(x *bonzai.Cmd, _ ...string) error {
			order = append(order, `install`)
			return nil
		}},
	}

	if err := InitCmd.Do(InitCmd); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(order) != 2 || order[0] != `setup` || order[1] != `install` {
		t.Errorf("call order = %v, want [setup install]", order)
	}
}
