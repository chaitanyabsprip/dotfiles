package cmdhelp

import (
	"strings"
	"testing"

	"github.com/rwxrob/bonzai"
)

func TestCommands(t *testing.T) {
	root := &bonzai.Cmd{
		Name: `root`,
		Cmds: []*bonzai.Cmd{
			{Name: `aa`, Short: `short a`},
			{Name: `bb`, Short: `short b`},
			{Name: `group`, Short: `has children`, Cmds: []*bonzai.Cmd{
				{Name: `child`, Short: `leaf`},
			}},
		},
	}

	out := Commands(root)

	for _, want := range []string{`aa - short a`, `bb - short b`, `group - has children`, `child - leaf`} {
		if !strings.Contains(out, want) {
			t.Fatalf("output missing %q, got:\n%s", want, out)
		}
	}

	wide := packGrid([]string{`a`, `b`, `c`, `d`}, 2)
	if len(wide) != 4 {
		t.Fatalf("expected 1 column per row when width is tight, got %v", wide)
	}

	narrow := packGrid([]string{`a`, `b`, `c`, `d`}, 80)
	if len(narrow) != 1 {
		t.Fatalf("expected a single row when width is plenty, got %v", narrow)
	}
}
