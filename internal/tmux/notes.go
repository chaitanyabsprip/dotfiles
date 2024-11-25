package tmux

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"
	"github.com/rwxrob/bonzai/fn/maps"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/pkg/fzf"
)

var NotesCmd = &bonzai.Cmd{
	Name:  `notes`,
	Alias: `n`,
	Short: `fuzzy search notes`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			args = append(args, ``)
		}
		return Notes(args[0], os.Getenv(`NOTESPATH`))
	},
}

func Notes(name, root string) error {
	files := getFiles(root)
	if len(name) == 0 {
		name = selectNote(files)
	}
	if len(name) == 0 {
		return nil
	}
	matches := filt.HasSuffix(files, name)
	if len(matches) == 0 {
		return nil
	}
	name = matches[0]
	exe := run.ExeName()
	err := run.Exec(exe, `tmux`, `x`, `sz`, `notes`)
	if err != nil {
		return err
	}
	err = run.Exec(
		`tmux`, `send-keys`,
		`-t`, `notes`,
		`cd`, root, `Enter`,
	)
	if err != nil {
		return err
	}
	return run.Exec(
		`tmux`, `send-keys`,
		`-t`, `notes`,
		`nvim `, name, `Enter`,
	)
}

func selectNote(files []string) string {
	out, err := fzf.Select(maps.Base(files), `--tmux`)
	if err != nil {
		return ``
	}
	return out
}

func getFiles(root string) []string {
	result := make([]string, 0)
	err := filepath.WalkDir(
		root,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if filepath.Ext(path) != ".md" {
				return nil
			}
			result = append(result, path)
			return nil
		},
	)
	if err != nil {
		return []string{}
	}
	return result
}
