package dot

import (
	"testing"
)

func TestDepsCmd_NoArgs_ReturnsError(t *testing.T) {
	err := DepsCmd.Do(DepsCmd)
	if err == nil {
		t.Error("expected error when called with no args, got nil")
	}
}

func TestDepsCmd_WithTool_PrintsNoDepsMessage(t *testing.T) {
	err := DepsCmd.Do(DepsCmd, "sometool")
	if err != nil {
		t.Errorf("unexpected error for tool arg: %v", err)
	}
}
