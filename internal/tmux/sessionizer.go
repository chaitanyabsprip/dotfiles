package tmux

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/fn/filt"

	"github.com/Chaitanyabsprip/dotfiles/pkg/fzf"
	"github.com/Chaitanyabsprip/dotfiles/pkg/tmux"
	"github.com/Chaitanyabsprip/dotfiles/x/depends"
	"github.com/Chaitanyabsprip/dotfiles/x/workdirs"
)

var SessionizerCmd = &bonzai.Cmd{
	Name:  `sessionizer`,
	Alias: `sz`,
	Short: `create a new or switch to an existing session`,
	Init: func(x *bonzai.Cmd, args ...string) error {
		depends.On(nil, `fzf`, `tmux`)
		return nil
	},
	Do: func(x *bonzai.Cmd, args ...string) error {
		// should always exit with non-zero status
		if len(args) > 0 {
			Sessionizer(args[0])
			return nil
		}
		Sessionizer(``)
		return nil
	},
}

func Sessionizer(path string) error {
	newPath := path
	if len(newPath) == 0 {
		newPath = selectPath()
	}
	if len(newPath) == 0 {
		return fmt.Errorf(`no path selected`)
	}
	sessionName := strings.ReplaceAll(
		filepath.Base(newPath),
		`.`,
		`_`,
	)
	if !tmux.IsActive() {
		return tmux.NewSession(
			tmux.Session{Name: sessionName, Path: newPath},
		)
	}

	fopts := tmux.Session{Path: newPath}
	if oldName, _ := tmux.FindSession(fopts); len(oldName) > 0 {
		return tmux.SwitchClient(oldName)
	}
	if _, oldPath := tmux.FindSession(tmux.Session{Name: sessionName}); len(
		oldPath,
	) > 0 {

		sessionName, oldSessionNewName := diffPath(newPath, oldPath)
		sessionName = strings.ReplaceAll(
			filepath.Base(newPath),
			`.`,
			`_`,
		)
		tmux.RenameSession(sessionName, oldSessionNewName)
	}
	err := tmux.NewSession(
		tmux.Session{Name: sessionName, Path: newPath},
	)
	if err != nil {
		return err
	}
	return tmux.SwitchClient(sessionName)
}

func selectPath() string {
	dirs := append(workdirs.Workdirs(), workdirs.Worktrees()...)
	out, err := fzf.Select(
		workdirs.Shorten(dirs),
		`--tmux`, `45%`,
		`--border`,
		`--border-label`, ` Sessionizer `,
		`--border-label-pos`, `6:bottom`,
	)
	if err != nil || len(out) == 0 {
		return ``
	}
	matches := filt.HasSuffix(dirs, out)
	if len(matches) == 0 {
		return ``
	} else {
		return matches[len(matches)-1]
	}
}

func diffPath(p1, p2 string) (string, string) {
	name1 := filepath.Base(p1)
	name2 := filepath.Base(p2)

	for name1 == name2 && len(name1) > 0 && len(name2) > 0 {
		p1 = filepath.Dir(p1)
		p2 = filepath.Dir(p2)
		name1 = filepath.Join(
			filepath.Base(p1),
			filepath.Base(name1),
		)
		name2 = filepath.Join(
			filepath.Base(p2),
			filepath.Base(name2),
		)
	}
	return name1, name2
}
