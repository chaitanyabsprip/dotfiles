package workdirs

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// worktree is one entry of `git worktree list --porcelain`.
type worktree struct {
	Path   string
	Branch string // short name; empty when detached or bare
	Bare   bool
}

// project describes where a repository keeps its worktrees.
type project struct {
	Root   string // where worktrees are created
	Repo   string // a directory git commands can run in
	Common string // absolute git common dir (.git or .bare)
}

// git runs git in dir and returns its trimmed stdout. On failure the
// error carries git's stderr.
func git(dir string, args ...string) (string, error) {
	cmd := exec.Command(`git`, append([]string{`-C`, dir}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return ``, fmt.Errorf(
			`git %s: %w: %s`,
			strings.Join(args, ` `), err, strings.TrimSpace(stderr.String()),
		)
	}
	return strings.TrimSpace(string(out)), nil
}

func parseWorktrees(porcelain string) []worktree {
	var out []worktree
	for _, line := range strings.Split(porcelain, "\n") {
		if p, ok := strings.CutPrefix(line, `worktree `); ok {
			out = append(out, worktree{Path: p})
			continue
		}
		if len(out) == 0 {
			continue
		}
		cur := &out[len(out)-1]
		if b, ok := strings.CutPrefix(line, `branch `); ok {
			cur.Branch = strings.TrimPrefix(b, `refs/heads/`)
		} else if line == `bare` {
			cur.Bare = true
		}
	}
	return out
}

func listWorktrees(dir string) ([]worktree, error) {
	out, err := git(dir, `worktree`, `list`, `--porcelain`)
	if err != nil {
		return nil, err
	}
	return parseWorktrees(out), nil
}

// openProject finds the project containing dir. dir may be any checkout
// of the project, its .bare dir, or the project dir itself.
func openProject(dir string) (project, error) {
	repo, err := repoDir(dir)
	if err != nil {
		return project{}, err
	}
	common, err := git(repo, `rev-parse`, `--path-format=absolute`, `--git-common-dir`)
	if err != nil {
		return project{}, err
	}
	root := filepath.Dir(common) // project/.bare -> project, repo/.git -> repo
	if filepath.Base(common) == `.git` {
		if b := filepath.Base(root); b == `root` || b == `main` {
			root = filepath.Dir(root) // home-base checkout
		}
	}
	return project{Root: root, Repo: repo, Common: common}, nil
}

// repoDir returns dir when it is inside a repository, else the project's
// .bare, root or main directory.
func repoDir(dir string) (string, error) {
	if _, err := git(dir, `rev-parse`, `--git-dir`); err == nil {
		return dir, nil
	}
	for _, name := range []string{`.bare`, `root`, `main`} {
		p := filepath.Join(dir, name)
		if _, err := git(p, `rev-parse`, `--git-dir`); err == nil {
			return p, nil
		}
	}
	return ``, fmt.Errorf(`%s is not inside a git project`, dir)
}

// samePath reports whether a and b name the same file, resolving
// symlinks (macOS temp dirs sit behind /var -> /private/var).
func samePath(a, b string) bool { return resolve(a) == resolve(b) }

func resolve(p string) string {
	if r, err := filepath.EvalSymlinks(p); err == nil {
		return r
	}
	return filepath.Clean(p)
}
