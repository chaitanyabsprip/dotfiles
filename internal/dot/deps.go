// Package dot provides dependency management infrastructure for dotfile tools.
package dot

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	"github.com/Chaitanyabsprip/dotfiles/internal/bash"
	"github.com/Chaitanyabsprip/dotfiles/internal/ohmyposh"
	"github.com/Chaitanyabsprip/dotfiles/internal/tmux"
	"github.com/Chaitanyabsprip/dotfiles/x/depends"
)

// toolDeps holds the tools whose runtime dependencies go beyond
// themselves, grounded in what their install/setup code actually shells
// out to. A tool absent here (most of them) just renders as a leaf —
// nothing to declare doesn't mean "unknown", per ADR-0011.
var toolDeps = map[string][]depends.Dep{
	`bash`:     bash.Deps,
	`ohmyposh`: ohmyposh.Deps,
	`tmux`:     tmux.Deps,
}

// DepsCmd shows a tool's dependency tree, per ADR-0011.
var DepsCmd = &bonzai.Cmd{
	Name:  `deps`,
	Short: `show dependency tree for a tool`,
	Comp:  comp.Cmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			return fmt.Errorf("usage: deps <tool> — specify a tool to check its dependencies")
		}
		depends.PrintTree(args[0], toolDeps[args[0]])
		return nil
	},
}
