// Package cmdhelp renders a bonzai command's subcommand list packed into
// columns across the terminal width, instead of bonzai's default one-per-line
// listing.
package cmdhelp

import (
	"fmt"
	"os"
	"strings"

	"github.com/rwxrob/bonzai"
	"golang.org/x/term"
)

const maxDepth = 2

// Commands is a drop-in replacement for the bonzai mark package's
// "commands" template func, assignable via bonzai.Cmd.Funcs.
func Commands(x *bonzai.Cmd) string {
	lines := renderChildren(x, 0, termWidth())
	return strings.Join(lines, "\n") + "\n"
}

func renderChildren(parent *bonzai.Cmd, depth, width int) []string {
	if depth >= maxDepth {
		return nil
	}
	children := visible(parent.Cmds)
	if len(children) == 0 {
		return nil
	}
	indent := strings.Repeat("  ", depth)
	avail := max(width-len(indent), 10)

	// Cells for leaf siblings share one column width across the whole
	// list, even when branch siblings split them into separate runs, so
	// every run of columns lines up with the others.
	cells := make([]string, len(children))
	leaf := make([]bool, len(children))
	widest := 0
	for i, c := range children {
		cell := c.Name
		if c.Short != "" {
			cell += " - " + c.Short
		}
		cells[i] = cell
		leaf[i] = len(visible(c.Cmds)) == 0
		if leaf[i] {
			if l := len([]rune(cell)); l > widest {
				widest = l
			}
		}
	}

	var out, run []string
	flush := func() {
		if len(run) == 0 {
			return
		}
		for _, line := range packGrid(run, avail, widest) {
			out = append(out, indent+line)
		}
		run = nil
	}
	for i, c := range children {
		if leaf[i] {
			run = append(run, cells[i])
			continue
		}
		flush()
		out = append(out, indent+cells[i])
		out = append(out, renderChildren(c, depth+1, width)...)
	}
	flush()
	return out
}

// packGrid lays cells out left-to-right, top-to-bottom in as many
// columns of the given width as fit within width.
func packGrid(cells []string, width, widest int) []string {
	if len(cells) == 0 {
		return nil
	}
	const gap = 2
	cols := min(max((width+gap)/(widest+gap), 1), len(cells))
	rows := (len(cells) + cols - 1) / cols
	lines := make([]string, 0, rows)
	for r := range rows {
		var b strings.Builder
		for c := range cols {
			i := r*cols + c
			if i >= len(cells) {
				break
			}
			if c == cols-1 || i == len(cells)-1 {
				b.WriteString(cells[i])
			} else {
				fmt.Fprintf(&b, "%-*s", widest+gap, cells[i])
			}
		}
		lines = append(lines, b.String())
	}
	return lines
}

func visible(cmds []*bonzai.Cmd) []*bonzai.Cmd {
	var out []*bonzai.Cmd
	for _, c := range cmds {
		if c != nil && !c.IsHidden() {
			out = append(out, c)
		}
	}
	return out
}

func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}
