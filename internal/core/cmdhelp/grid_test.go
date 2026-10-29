package cmdhelp

import (
	"fmt"
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

	wide := packGrid([]string{`a`, `b`, `c`, `d`}, 2, 1)
	if len(wide) != 4 {
		t.Fatalf("expected 1 column per row when width is tight, got %v", wide)
	}

	narrow := packGrid([]string{`a`, `b`, `c`, `d`}, 80, 1)
	if len(narrow) != 1 {
		t.Fatalf("expected a single row when width is plenty, got %v", narrow)
	}
}

func TestCommandsSharedColumnWidth(t *testing.T) {
	// The branch "group" splits the leaves into two runs (aa/cc and
	// dd/bb). Both runs must still share one column width - taken from
	// the longest leaf cell anywhere in the list - so their columns line
	// up instead of each run picking its own width.
	root := &bonzai.Cmd{
		Name: `root`,
		Cmds: []*bonzai.Cmd{
			{Name: `aa`, Short: `a`},
			{Name: `cc`, Short: `c`},
			{Name: `group`, Short: `has children`, Cmds: []*bonzai.Cmd{
				{Name: `child`, Short: `leaf`},
			}},
			{Name: `dd`, Short: `a much longer short description`},
			{Name: `bb`, Short: `b`},
		},
	}

	lines := renderChildren(root, 0, 200)
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines (aa/cc run, group, its child, dd/bb run), got %q", lines)
	}

	widest := len(`dd - a much longer short description`)
	wantFirstRow := fmt.Sprintf("%-*s%s", widest+2, `aa - a`, `cc - c`)
	if lines[0] != wantFirstRow {
		t.Fatalf("first run not padded to the shared column width:\ngot:  %q\nwant: %q", lines[0], wantFirstRow)
	}
	wantLastRow := fmt.Sprintf("%-*s%s", widest+2, `dd - a much longer short description`, `bb - b`)
	if lines[3] != wantLastRow {
		t.Fatalf("second run not padded to the shared column width:\ngot:  %q\nwant: %q", lines[3], wantLastRow)
	}
}
