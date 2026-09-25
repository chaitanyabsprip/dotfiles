package depends

import "fmt"

// Dep is one node in a tool's dependency tree: an external binary it
// needs on PATH (Name), optionally noting why (Note), and any
// dependencies of its own (Deps) — see ADR-0011.
type Dep struct {
	Name string
	Note string
	Deps []Dep
}

// PrintTree renders deps as a recursive tree, e.g.:
//
//	tmux
//	├── fzf (sessionizer)
//	└── tmux (binary)
func PrintTree(root string, deps []Dep) {
	fmt.Println(root)
	printChildren(deps, ``)
}

func printChildren(deps []Dep, prefix string) {
	for i, d := range deps {
		last := i == len(deps)-1
		branch, nextPrefix := `├── `, prefix+`│   `
		if last {
			branch, nextPrefix = `└── `, prefix+`    `
		}
		label := d.Name
		if d.Note != `` {
			label = fmt.Sprintf(`%s (%s)`, d.Name, d.Note)
		}
		fmt.Println(prefix + branch + label)
		printChildren(d.Deps, nextPrefix)
	}
}
