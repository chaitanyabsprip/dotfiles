package depends

import (
	"bytes"
	"os"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf(`os.Pipe: %v`, err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestPrintTreeMatchesADR0011Example(t *testing.T) {
	deps := []Dep{
		{Name: `fzf`, Note: `sessionizer`},
		{Name: `tmux`, Note: `binary`},
		{Name: `ohmyposh`, Deps: []Dep{{Name: `unzip`}}},
	}
	want := "tmux\n" +
		"├── fzf (sessionizer)\n" +
		"├── tmux (binary)\n" +
		"└── ohmyposh\n" +
		"    └── unzip\n"

	got := captureStdout(t, func() { PrintTree(`tmux`, deps) })
	if got != want {
		t.Fatalf("PrintTree output mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestPrintTreeLeaf(t *testing.T) {
	got := captureStdout(t, func() { PrintTree(`alacritty`, nil) })
	if got != "alacritty\n" {
		t.Fatalf(`PrintTree with no deps = %q, want %q`, got, "alacritty\n")
	}
}
