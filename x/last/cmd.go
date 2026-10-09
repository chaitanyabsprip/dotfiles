package last

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
	"github.com/Chaitanyabsprip/dotfiles/pkg/prompt"
)

// lookup finds a path to act on, given a directory to search.
type lookup func(dir string) (string, error)

var Cmd = &bonzai.Cmd{
	Name:    `last`,
	Short:   `find, move or copy the newest file or directory`,
	Usage:   `[path]`,
	MaxArgs: 1,
	Long: `
Prints the newest entry (file or directory) in [path]. [path] defaults
to ` + "`$DOWNLOADS`" + ` (or ` + "`~/downloads`" + `); pass ` + "`.`" + ` for the current
directory instead. See 'last help' for the dir/file/mv/cp/edit
commands.`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{dirCmd, fileCmd, mvCmd, cpCmd, editCmd, rmCmd, help.Cmd},
	Do:   doPrint(Find),
}

// downloadsDir is $DOWNLOADS, falling back to ~/downloads when unset.
func downloadsDir() string {
	if env.Downloads != `` {
		return env.Downloads
	}
	return filepath.Join(env.Home, `downloads`)
}

// pathArg returns args[0], or downloadsDir() when args is empty. Pass
// "." explicitly for the current directory.
func pathArg(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	return downloadsDir()
}

// doPrint returns a Do that prints what find locates in pathArg(args).
func doPrint(find lookup) func(_ *bonzai.Cmd, args ...string) error {
	return func(_ *bonzai.Cmd, args ...string) error {
		n, err := find(pathArg(args))
		if err != nil {
			return err
		}
		fmt.Println(n)
		return nil
	}
}

var dirCmd = &bonzai.Cmd{
	Name:    `dir`,
	Short:   `print the newest subdirectory in path`,
	Usage:   `[path]`,
	MaxArgs: 1,
	Long: `
Prints the newest subdirectory directly inside [path] (default:
` + "`$DOWNLOADS`" + `; pass ` + "`.`" + ` for the current directory).`,
	Do: doPrint(FindDir),
}

var fileCmd = &bonzai.Cmd{
	Name:    `file`,
	Short:   `print the newest file in path`,
	Usage:   `[path]`,
	MaxArgs: 1,
	Long: `
Prints the newest non-directory entry directly inside [path] (default:
` + "`$DOWNLOADS`" + `; pass ` + "`.`" + ` for the current directory).`,
	Do: doPrint(FindFile),
}

var editCmd = &bonzai.Cmd{
	Name:    `edit`,
	Alias:   `e`,
	Short:   `open the newest file in path`,
	Usage:   `[path]`,
	MaxArgs: 1,
	Long: `
Opens the newest non-directory entry directly inside [path] (default:
` + "`$DOWNLOADS`" + `; pass ` + "`.`" + ` for the current directory) in
` + "`$VISUAL`" + `, ` + "`$EDITOR`" + `, or nvim, whichever is set first.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		dir := pathArg(args)
		n, err := FindFile(dir)
		if err != nil {
			return err
		}
		if n == `` {
			return fmt.Errorf(`no files in %s`, dir)
		}
		return run.Exec(editorCmd(), n)
	},
}

// editorCmd picks $VISUAL, then $EDITOR, then nvim.
func editorCmd() string {
	if env.Visual != `` {
		return env.Visual
	}
	if env.Editor != `` {
		return env.Editor
	}
	return `nvim`
}

var rmCmd = &bonzai.Cmd{
	Name:    `rm`,
	Short:   `delete the newest download`,
	Usage:   `[path]`,
	MaxArgs: 1,
	Long: `
Deletes the newest entry (file or directory) in [path] (default:
` + "`$DOWNLOADS`" + `; pass ` + "`.`" + ` for the current directory).
Irreversible, so it asks twice: a y/N prompt, then a random word you
must type back exactly.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		dir := pathArg(args)
		n, err := Find(dir)
		if err != nil {
			return err
		}
		if n == `` {
			return fmt.Errorf(`no matching entry in %s`, dir)
		}
		fmt.Printf("About to delete: %s\n", n)
		in := bufio.NewReader(os.Stdin)
		if !prompt.Confirm(in, os.Stdout, `Are you sure?`) {
			return errors.New(`aborted`)
		}
		if !prompt.ConfirmWord(in, os.Stdout) {
			return errors.New(`aborted`)
		}
		return os.RemoveAll(n)
	},
}

var mvCmd = &bonzai.Cmd{
	Name:    `mv`,
	Short:   `move the newest download into pwd`,
	Usage:   `<path>`,
	NumArgs: 1,
	Cmds:    []*bonzai.Cmd{mvFileCmd, mvDirCmd, help.Cmd},
	Long: `
Moves the newest entry (file or directory) in ` + "`$DOWNLOADS`" + ` (or
` + "`~/downloads`" + `) to ` + "`<path>`" + ` in the current directory. See 'mv file'
and 'mv dir' to move only one type.`,
	Do: doTransfer(Find, moveOp),
}

var mvFileCmd = &bonzai.Cmd{
	Name:    `file`,
	Short:   `move the newest downloaded file into pwd`,
	Usage:   `<path>`,
	NumArgs: 1,
	Long: `
Moves the newest file in ` + "`$DOWNLOADS`" + ` (or ` + "`~/downloads`" + `) to
` + "`<path>`" + `, skipping any newer directory there.`,
	Do: doTransfer(FindFile, moveOp),
}

var mvDirCmd = &bonzai.Cmd{
	Name:    `dir`,
	Short:   `move the newest downloaded directory into pwd`,
	Usage:   `<path>`,
	NumArgs: 1,
	Long: `
Moves the newest subdirectory in ` + "`$DOWNLOADS`" + ` (or ` + "`~/downloads`" + `) to
` + "`<path>`" + `, skipping any newer file there.`,
	Do: doTransfer(FindDir, moveOp),
}

var cpCmd = &bonzai.Cmd{
	Name:    `cp`,
	Short:   `copy the newest download into pwd`,
	Usage:   `<path>`,
	NumArgs: 1,
	Cmds:    []*bonzai.Cmd{cpFileCmd, cpDirCmd, help.Cmd},
	Long: `
Copies the newest entry (file or directory) in ` + "`$DOWNLOADS`" + ` (or
` + "`~/downloads`" + `) to ` + "`<path>`" + ` in the current directory, leaving the
original in place. See 'cp file' and 'cp dir' to copy only one type.`,
	Do: doTransfer(Find, copyOp),
}

var cpFileCmd = &bonzai.Cmd{
	Name:    `file`,
	Short:   `copy the newest downloaded file into pwd`,
	Usage:   `<path>`,
	NumArgs: 1,
	Long: `
Copies the newest file in ` + "`$DOWNLOADS`" + ` (or ` + "`~/downloads`" + `) to
` + "`<path>`" + `, skipping any newer directory there.`,
	Do: doTransfer(FindFile, copyOp),
}

var cpDirCmd = &bonzai.Cmd{
	Name:    `dir`,
	Short:   `copy the newest downloaded directory into pwd`,
	Usage:   `<path>`,
	NumArgs: 1,
	Long: `
Copies the newest subdirectory in ` + "`$DOWNLOADS`" + ` (or ` + "`~/downloads`" + `) to
` + "`<path>`" + `, skipping any newer file there.`,
	Do: doTransfer(FindDir, copyOp),
}

// doTransfer returns a Do that finds a path with find in downloadsDir(),
// hands it to op along with args[0] resolved against the working
// directory, and prints the result.
func doTransfer(find lookup, op func(from, to string) error) func(_ *bonzai.Cmd, args ...string) error {
	return func(_ *bonzai.Cmd, args ...string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		from, to, err := transfer(find, wd, args[0], op)
		if err != nil {
			return err
		}
		fmt.Printf("From: %s\nTo:   %s\n", from, to)
		return nil
	}
}

// transfer finds a path with find in downloadsDir(), resolves dest
// against wd, and hands both to op (a move or a copy).
func transfer(find lookup, wd, dest string, op func(from, to string) error) (from, to string, err error) {
	dir := downloadsDir()
	from, err = find(dir)
	if err != nil {
		return ``, ``, err
	}
	if from == `` {
		return ``, ``, fmt.Errorf(`no matching entry in %s`, dir)
	}
	to = filepath.Join(wd, dest)
	if err := op(from, to); err != nil {
		return ``, ``, err
	}
	return from, to, nil
}

func moveOp(from, to string) error { return run.Exec(`mv`, from, to) }

// copyOp copies from to to, recursing when from is a directory.
func copyOp(from, to string) error {
	info, err := os.Stat(from)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return run.Exec(`cp`, `-r`, from, to)
	}
	return run.Exec(`cp`, from, to)
}
