package dot

import (
	"testing"

	"github.com/rwxrob/bonzai"
)

func TestEditCmdDoNoArgsReturnsError(t *testing.T) {
	if err := EditCmd.Do(EditCmd); err == nil {
		t.Error("EditCmd.Do() with no args should return error, got nil")
	}
}

func TestEditCmdDoSpecificTool(t *testing.T) {
	orig := EditCmds
	defer func() { EditCmds = orig }()

	called := false
	EditCmds = []*bonzai.Cmd{
		{Name: `mock`, Do: func(x *bonzai.Cmd, _ ...string) error { called = true; return nil }},
	}

	if err := EditCmd.Do(EditCmd, `mock`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("mock edit was not called for Do('mock')")
	}
}
