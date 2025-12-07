// Package tmux
package tmux

import (
	"fmt"
	"path/filepath"
	"sort"
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
	sessionName := resolveSessionName(newPath)

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

		oldSessionNewName, sessionName := reconcileSessionName(
			newPath,
			oldPath,
		)
		sessionName = strings.ReplaceAll(sessionName, `.`, `_`)
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

func resolveSessionName(path string) string {
	base := filepath.Base(path)
	if base == "root" {
		base = fmt.Sprintf("%s/r", filepath.Base(filepath.Dir(path)))
	}
	return strings.ReplaceAll(base, `.`, `_`)
}

func selectPath() string {
	dirs := workdirs.Workdirs()
	trees := workdirs.Worktrees()
	dirs = append(dirs, trees...)
	sort.Strings(dirs)
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

func reconcileSessionName(oldPath, newPath string) (string, string) {
	dirA, baseA := filepath.Split(oldPath)
	dirB, baseB := filepath.Split(newPath)
	if oldPath == newPath {
		if baseA == "root" {
			parentBase := filepath.Base(filepath.Dir(oldPath))
			return fmt.Sprintf(
					"old-%s/%s",
					parentBase,
					baseA,
				), fmt.Sprintf(
					"%s/%s",
					parentBase,
					baseB,
				)
		}
		return fmt.Sprint("old-", baseA), baseB
	}
	fmt.Println(dirA, baseA, dirB, baseB)
	if baseA == baseB {
		nBaseB := baseB
		nBaseA := baseA
		nnBaseB := nBaseB
		for nnBaseB == baseA {
			dirA, nBaseA = filepath.Split(strings.TrimRight(dirA, "/"))
			baseA = filepath.Join(nBaseA, baseA)
			dirB, nBaseB = filepath.Split(strings.TrimRight(dirB, "/"))
			nnBaseB = filepath.Join(nBaseB, nnBaseB)
		}
		if baseB == "root" {
			return baseA, nnBaseB
		}
		return baseA, baseB
	}
	return baseA, baseB
}
