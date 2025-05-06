package workdirs

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

// Cmd provides access to the `work` command suite, which lists work
// directories and Git worktree repositories. Use SHORT=1 for compact
// output.
var Cmd = &bonzai.Cmd{
	Name:  `work`,
	Short: `list work directories and Git worktrees`,
	Long: `The 'work' command lists local work directories and
Git worktree repositories. It supports compact output when SHORT=1 is
set.`,
	Def:  helpCmd,
	Cmds: []*bonzai.Cmd{helpCmd, dirsCmd, treeCmd},
	Comp: comp.Cmds,
}

// helpCmd shows detailed help for each subcommand or a summary of all
// commands.
var helpCmd = &bonzai.Cmd{
	Name:  `help`,
	Short: `show help for a command`,
	Opts:  `dirs|trees`,
	Comp:  comp.Opts,
	Long: `
Displays usage and description for all available commands, or detailed
information about a specified command if provided.

USAGE
  work help [command]

EXAMPLES
  work help       # shows usage for all commands
  work help dirs  # shows help for 'dirs' command
  work help trees # shows help for 'trees' command`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		cmds := []*bonzai.Cmd{dirsCmd, treeCmd}
		if len(args) == 0 {
			fmt.Printf("%s - %s\n\n", x.Name, x.Short)
			fmt.Println("COMMANDS:")
			for _, cmd := range cmds {
				fmt.Printf("  %-10s - %s\n", cmd.Name, cmd.Short)
			}
			fmt.Printf(
				"\nUse 'help [command]' for more information about a command.\n",
			)
		} else if len(args) == 1 {
			var matched bool
			for _, cmd := range cmds {
				if cmd.Name == args[0] {
					matched = true
					fmt.Printf("%s - %s\n\n", cmd.Name, cmd.Short)
					fmt.Println("DESCRIPTION:")
					fmt.Printf("  %s\n\n", cmd.Long)
				}
			}
			if !matched {
				return fmt.Errorf("unknown command %q", args[0])
			}
		} else {
			return fmt.Errorf("too many arguments")
		}
		return nil
	},
}

// dirsCmd lists work directories with optional compact output.
var dirsCmd = &bonzai.Cmd{
	Name:  `dirs`,
	Alias: `d`,
	Short: `list work directories`,
	Long: `
Lists local work directories, including Git worktree repositories.
Set SHORT=1 for a compact output format.

ENVIRONMENT VARIABLES
  SHORT|W_SHORT     Set SHORT=1 to display compact output
`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		short := len(os.Getenv(`SHORT`)) > 0 ||
			len(os.Getenv(`W_SHORT`)) > 0
		var out string
		dirs := Workdirs()
		trees := Worktrees()
		dirs = append(dirs, trees...)
		sort.Strings(dirs)
		if !short {
			out = strings.Join(dirs, "\n")
		} else {
			out = strings.Join(Shorten(dirs), "\n")
		}
		fmt.Println(out)
		return nil
	},
}

// treeCmd lists Git worktrees with optional compact output.
var treeCmd = &bonzai.Cmd{
	Name:  `trees`,
	Alias: `t`,
	Short: `list Git worktrees`,
	Long: `
Lists Git worktrees within the specified directories.
Use SHORT=1 for a compact output format.

ENVIRONMENT VARIABLES
  SHORT     Set SHORT=1 to display compact output
  W_SHORT   Alternate way to set SHORT=1 for compact output
`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		short := len(os.Getenv(`SHORT`)) > 0 ||
			len(os.Getenv(`W_SHORT`)) > 0
		var out string
		treeList := Worktrees()
		sort.Strings(treeList)
		if !short {
			out = strings.Join(treeList, "\n")
		} else {
			out = strings.Join(
				Shorten(treeList),
				"\n",
			)
		}
		fmt.Println(out)
		return nil
	},
}
