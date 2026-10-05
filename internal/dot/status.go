package dot

import (
	"fmt"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"

	e "github.com/Chaitanyabsprip/dotfiles/internal/core/embed"
)

// StatusCmd lists every drifted or never-deployed config file across all
// tools (or one tool), without touching disk or printing diffs — a
// preview before `dot setup` would decide, per file, to skip-and-warn,
// create, or update.
var StatusCmd = &bonzai.Cmd{
	Name:  `status`,
	Short: `list drifted and un-hydrated config files`,
	Comp:  comp.Cmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		tool := ``
		if len(args) > 0 && args[0] != `all` {
			tool = args[0]
		}
		results, err := checkDrift(tool, e.CheckList)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			fmt.Println(`clean — nothing drifted or un-hydrated`)
			return nil
		}
		for _, cmd := range driftCheckCmds {
			result, ok := results[cmd.Name]
			if !ok {
				continue
			}
			for _, path := range result.Drifted {
				fmt.Printf("%s\t%s\tdrifted\n", cmd.Name, path)
			}
			for _, path := range result.Missing {
				fmt.Printf("%s\t%s\tmissing\n", cmd.Name, path)
			}
		}
		return nil
	},
}

// DiffCmd shows the diff for every drifted file in a tool (or all tools),
// and for every never-deployed file, the diff against nothing (i.e. what
// `dot setup` would create).
var DiffCmd = &bonzai.Cmd{
	Name:  `diff`,
	Short: `show diff for a tool's drifted and un-hydrated config files`,
	Comp:  comp.Cmds,
	Do: func(x *bonzai.Cmd, args ...string) error {
		tool := ``
		if len(args) > 0 && args[0] != `all` {
			tool = args[0]
		}
		results, err := checkDrift(tool, e.CheckDiff)
		if err != nil {
			return err
		}
		if len(results) == 0 {
			fmt.Println(`clean — nothing drifted or un-hydrated`)
		}
		return nil
	},
}

// toolCheck is one tool's findings from a checkDrift run.
type toolCheck struct {
	Drifted []string
	Missing []string
}

// checkDrift runs the named tool's setup (or every driftCheckCmds tool, if
// name is "") in the given e.CheckMode and returns each tool alongside
// what it found. Supports StatusCmd and DiffCmd above.
func checkDrift(name string, mode e.CheckMode) (map[string]toolCheck, error) {
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

	result := map[string]toolCheck{}
	for _, cmd := range cmds {
		e.Drifted = nil
		e.Missing = nil
		if err := cmd.Run(``); err != nil {
			return nil, err
		}
		if len(e.Drifted) > 0 || len(e.Missing) > 0 {
			result[cmd.Name] = toolCheck{Drifted: e.Drifted, Missing: e.Missing}
		}
	}
	return result, nil
}

// nonDriftCheckable are SetupCmds that don't go through
// internal/core/embed, so they have no manifest/drift concept to check
// and running them here would perform a real side effect instead of a
// preview: claude's setup merges ~/.claude/settings.json directly;
// neovim's clones/pulls its config from a separate git repo.
var nonDriftCheckable = map[string]bool{
	`claude`: true,
	`neovim`: true,
}

// driftCheckCmds is SetupCmds minus nonDriftCheckable.
var driftCheckCmds = func() []*bonzai.Cmd {
	cmds := make([]*bonzai.Cmd, 0, len(SetupCmds))
	for _, cmd := range SetupCmds {
		if !nonDriftCheckable[cmd.Name] {
			cmds = append(cmds, cmd)
		}
	}
	return cmds
}()
