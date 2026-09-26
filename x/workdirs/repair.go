package workdirs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai"
)

var repairCmd = &bonzai.Cmd{
	Name:    `repair`,
	Alias:   `r`,
	Short:   `reconnect worktrees after moving a project`,
	Usage:   `repair [dir]`,
	MaxArgs: 1,
	Long: `
Reconnects a project's worktrees after the project directory was moved
or renamed. Run it from the project's new location (or pass that
directory). It finds the main repository and every linked worktree up to
three levels down and runs 'git worktree repair' on them, then prints
the worktrees it repaired.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		dir := `.`
		if len(args) == 1 {
			dir = args[0]
		}
		abs, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		linked, err := repairProject(abs)
		if err != nil {
			return err
		}
		for _, p := range linked {
			fmt.Println(p)
		}
		return nil
	},
}

// repairProject reconnects every linked worktree under root to its
// repository and returns the worktrees it handed to git.
func repairProject(root string) ([]string, error) {
	repo, linked, err := scanProject(root)
	if err != nil || len(linked) == 0 {
		return nil, err
	}
	if _, err := git(repo, append([]string{`worktree`, `repair`}, linked...)...); err != nil {
		return nil, err
	}
	return linked, nil
}

// scanProject finds, up to three levels under root, the shallowest main
// repository (a .bare dir, or the dir holding a .git dir) and every
// linked worktree (a dir whose .git file points into worktrees/).
func scanProject(root string) (repo string, linked []string, err error) {
	repoDepth := -1
	err = filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		depth := strings.Count(strings.TrimPrefix(p, root), string(filepath.Separator))
		if d.IsDir() && (depth > 3 || d.Name() == `node_modules`) {
			return fs.SkipDir
		}
		if d.Name() != `.git` && d.Name() != `.bare` {
			return nil
		}
		if d.IsDir() {
			found := p // a .bare dir is the repository itself
			if d.Name() == `.git` {
				found = filepath.Dir(p)
			}
			if repoDepth < 0 || depth < repoDepth {
				repo, repoDepth = found, depth
			}
			return fs.SkipDir
		}
		if d.Name() == `.git` && isWorktreeGitfile(p) {
			linked = append(linked, filepath.Dir(p))
		}
		return nil
	})
	if err == nil && repo == `` {
		err = fmt.Errorf(`no git repository under %s`, root)
	}
	return repo, linked, err
}

// isWorktreeGitfile reports whether the .git file at p belongs to a
// linked worktree. Submodules use .git files too, but they point into
// modules/, not worktrees/.
func isWorktreeGitfile(p string) bool {
	b, err := os.ReadFile(p)
	return err == nil && strings.Contains(filepath.ToSlash(string(b)), `/worktrees/`)
}
