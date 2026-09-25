package dot

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
)

// StatusCmd lists every drifted file across all tools (or one tool),
// without touching disk or printing diffs — a preview before `dot setup`
// would decide, per file, to skip-and-warn or deploy.
var StatusCmd = &bonzai.Cmd{
	Name:  `status`,
	Short: `list drifted config files`,
	Comp:  comp.Cmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		tool := ``
		if len(args) > 0 && args[0] != `all` {
			tool = args[0]
		}
		drift, err := checkDrift(tool, e.CheckList)
		if err != nil {
			return err
		}
		if len(drift) == 0 {
			fmt.Println(`clean — no drifted config files`)
			return nil
		}
		for _, cmd := range driftCheckCmds {
			paths, ok := drift[cmd.Name]
			if !ok {
				continue
			}
			for _, path := range paths {
				fmt.Printf("%s\t%s\n", cmd.Name, path)
			}
		}
		return nil
	},
}

// DiffCmd shows the diff for every drifted file in a tool (or all tools).
var DiffCmd = &bonzai.Cmd{
	Name:  `diff`,
	Short: `show diff for a tool's drifted config files`,
	Comp:  comp.Cmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		tool := ``
		if len(args) > 0 && args[0] != `all` {
			tool = args[0]
		}
		drift, err := checkDrift(tool, e.CheckDiff)
		if err != nil {
			return err
		}
		if len(drift) == 0 {
			fmt.Println(`clean — no drifted config files`)
		}
		return nil
	},
}

// checkDrift runs the named tool's setup (or every driftCheckCmds tool, if
// name is "") in the given e.CheckMode and returns each tool alongside the
// dest paths it found drifted. Supports StatusCmd and DiffCmd above.
func checkDrift(name string, mode e.CheckMode) (map[string][]string, error) {
	cmds := driftCheckCmds
	if name != `` {
		found := false
		for _, cmd := range driftCheckCmds {
			if cmd.Name == name {
				cmds = []*bonzai.Cmd{cmd}
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf(`unknown tool: %s`, name)
		}
	}

	e.Check = mode
	defer func() { e.Check = e.CheckOff }()

	result := map[string][]string{}
	for _, cmd := range cmds {
		e.Drifted = nil
		if err := cmd.Run(``); err != nil {
			return nil, err
		}
		if len(e.Drifted) > 0 {
			result[cmd.Name] = e.Drifted
		}
	}
	return result, nil
}

// driftCheckCmds is SetupCmds minus claude.SetupCmd: claude's setup merges
// ~/.claude/settings.json directly and never goes through
// internal/core/embed, so it has no manifest/drift concept to check and
// running it here would perform a real write instead of a preview.
var driftCheckCmds = func() []*bonzai.Cmd {
	cmds := make([]*bonzai.Cmd, 0, len(SetupCmds)-1)
	for _, cmd := range SetupCmds {
		if cmd.Name != `claude` {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}()
