package workdirs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/bonzai"
)

var addCmd = &bonzai.Cmd{
	Name:    `add`,
	Alias:   `a`,
	Short:   `create a worktree for a branch`,
	Usage:   `add <branch>`,
	NumArgs: 1,
	Long: `
Creates a worktree for <branch> under the project root and prints its
path, so 'cd "$(work add feat/x)"' lands in it. The path keeps the
branch's first and last segments: feat/mod/TICKET-1/short-desc becomes
feat/short-desc. An existing local or remote branch is checked out; a
new name creates the branch. If the branch already has a worktree, that
worktree's path is printed instead.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		path, err := addWorktree(wd, args[0])
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	},
}

// worktreePathFor maps a branch to its worktree dir relative to the
// project root: "feat/mod/TICKET-1/short-desc" -> "feat/short-desc".
func worktreePathFor(branch string) string {
	parts := strings.Split(strings.Trim(branch, `/`), `/`)
	if len(parts) == 1 {
		return parts[0]
	}
	return filepath.Join(parts[0], parts[len(parts)-1])
}

// addWorktree creates (or finds) the worktree for branch in the project
// containing dir and returns its path.
func addWorktree(dir, branch string) (string, error) {
	p, err := openProject(dir)
	if err != nil {
		return ``, err
	}
	wts, err := listWorktrees(p.Repo)
	if err != nil {
		return ``, err
	}
	for _, wt := range wts {
		if wt.Branch == branch {
			fmt.Fprintf(os.Stderr, "%s is already checked out at %s\n", branch, wt.Path)
			return wt.Path, nil
		}
	}
	path := filepath.Join(p.Root, worktreePathFor(branch))
	args := []string{`worktree`, `add`, `-b`, branch, path}
	if branchExists(p.Repo, branch) {
		// For a remote-only branch git creates a local tracking branch.
		args = []string{`worktree`, `add`, path, branch}
	}
	if _, err := git(p.Repo, args...); err != nil {
		return ``, err
	}
	return path, nil
}

// branchExists reports whether branch exists locally or on any remote.
func branchExists(repo, branch string) bool {
	if _, err := git(repo, `show-ref`, `--verify`, `--quiet`, `refs/heads/`+branch); err == nil {
		return true
	}
	out, _ := git(repo, `for-each-ref`, `--format=%(refname)`, `refs/remotes/*/`+branch)
	return out != ``
}
