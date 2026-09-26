package workdirs

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rwxrob/bonzai"
)

var prCmd = &bonzai.Cmd{
	Name:    `pr`,
	Short:   `check out a GitHub PR into a worktree`,
	Usage:   `pr <number> [path] | pr done [number] [path] | pr prune`,
	MinArgs: 1,
	MaxArgs: 2,
	Cmds:    []*bonzai.Cmd{prDoneCmd, prPruneCmd},
	Long: `
Checks out GitHub PR <number> into a worktree at pr/<number> (or [path])
under the project root and prints the worktree path, so
'cd "$(work pr 17)"' lands in it. Running it again for the same PR just
prints the path. If the PR's branch is already checked out elsewhere,
that worktree is printed instead. Needs the gh CLI.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		rel := ``
		if len(args) == 2 {
			rel = args[1]
		}
		path, err := prCheckout(wd, args[0], rel, os.Stdin, os.Stderr)
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	},
}

var prDoneCmd = &bonzai.Cmd{
	Name:    `done`,
	Alias:   `d`,
	Short:   `remove a PR worktree and its branch`,
	Usage:   `done [number] [path]`,
	MaxArgs: 2,
	Long: `
Removes the worktree for PR [number] and its local branch, then prints
the directory to be in afterwards: the project root when you were inside
the removed worktree. Without [number] the PR is taken from the current
pr/<number> worktree. A worktree with uncommitted or untracked files is
kept.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		var n, rel string
		if len(args) > 0 {
			n = args[0]
		}
		if len(args) > 1 {
			rel = args[1]
		}
		cd, err := prDone(wd, n, rel, os.Stderr)
		if err != nil {
			return err
		}
		fmt.Println(cd)
		return nil
	},
}

var prPruneCmd = &bonzai.Cmd{
	Name:   `prune`,
	Alias:  `p`,
	Short:  `remove worktrees of merged or closed PRs`,
	NoArgs: true,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		return prPrune(wd, os.Stdin, os.Stderr)
	},
}

// prInfo is the part of `gh pr view --json` that pr needs.
type prInfo struct {
	State       string `json:"state"`
	HeadRefName string `json:"headRefName"`
}

// Seams over the gh CLI; tests replace them.
var (
	ghPRView = func(dir, n string) (prInfo, error) {
		cmd := exec.Command(`gh`, `pr`, `view`, n, `--json`, `state,headRefName`)
		cmd.Dir = dir
		out, err := cmd.Output()
		if err != nil {
			return prInfo{}, err
		}
		var info prInfo
		return info, json.Unmarshal(out, &info)
	}
	ghPRCheckout = func(dir, n string) error {
		cmd := exec.Command(`gh`, `pr`, `checkout`, n)
		cmd.Dir, cmd.Stdout, cmd.Stderr = dir, os.Stderr, os.Stderr
		return cmd.Run()
	}
)

var (
	errAborted = errors.New(`aborted`)
	prNumber   = regexp.MustCompile(`^[0-9]+$`)
	prDir      = regexp.MustCompile(`^pr/([0-9]+)$`)
)

// prCheckout puts PR n in a worktree at rel (default pr/<n>) under the
// project root and returns its path.
func prCheckout(dir, n, rel string, in io.Reader, out io.Writer) (string, error) {
	if !prNumber.MatchString(n) {
		return ``, fmt.Errorf(`invalid PR number %q`, n)
	}
	p, err := openProject(dir)
	if err != nil {
		return ``, err
	}
	if rel == `` {
		rel = filepath.Join(`pr`, n)
	}
	path := filepath.Join(p.Root, rel)
	wts, err := listWorktrees(p.Repo)
	if err != nil {
		return ``, err
	}
	if wt, ok := findWorktree(wts, path); ok {
		fmt.Fprintf(out, "PR #%s is already checked out at %s\n", n, wt.Path)
		return wt.Path, nil
	}
	info, err := ghPRView(p.Repo, n)
	if err != nil {
		return ``, fmt.Errorf(`gh pr view %s: %w`, n, err)
	}
	if info.State == `MERGED` || info.State == `CLOSED` {
		q := fmt.Sprintf(`PR #%s is %s. Check it out anyway?`, n, info.State)
		if !confirm(in, out, q) {
			return ``, errAborted
		}
	}
	for _, wt := range wts {
		if info.HeadRefName != `` && wt.Branch == info.HeadRefName {
			fmt.Fprintf(out, "branch %s is already checked out at %s\n", wt.Branch, wt.Path)
			return wt.Path, nil
		}
	}
	// Detached first so git does not invent a branch named after the dir.
	if _, err := git(p.Repo, `worktree`, `add`, `--detach`, path); err != nil {
		return ``, err
	}
	if err := ghPRCheckout(path, n); err != nil {
		_, rmErr := git(p.Common, `worktree`, `remove`, `--force`, path)
		return ``, errors.Join(fmt.Errorf(`gh pr checkout %s: %w`, n, err), rmErr)
	}
	return path, nil
}

// prDone removes PR n's worktree (default pr/<n>) and its local branch.
// With n empty the PR is inferred from dir. It returns the directory the
// caller should be in afterwards.
func prDone(dir, n, rel string, out io.Writer) (string, error) {
	p, err := openProject(dir)
	if err != nil {
		return ``, err
	}
	if n == `` {
		if n = prFromDir(p.Root, dir); n == `` {
			return ``, fmt.Errorf(`no PR number given and %s is not in a pr/<number> worktree`, dir)
		}
	}
	if !prNumber.MatchString(n) {
		return ``, fmt.Errorf(`invalid PR number %q`, n)
	}
	if rel == `` {
		rel = filepath.Join(`pr`, n)
	}
	wts, err := listWorktrees(p.Repo)
	if err != nil {
		return ``, err
	}
	wt, ok := findWorktree(wts, filepath.Join(p.Root, rel))
	if !ok {
		fmt.Fprintf(out, "no worktree at %s, nothing to remove\n", rel)
		return dir, nil
	}
	// Decide before removal: once the dir is gone its symlinks can't resolve.
	inside := isInside(dir, wt.Path)
	info, _ := ghPRView(p.Repo, n)
	// Run from the common dir: dir may be the worktree being removed.
	if err := removeWorktree(p.Common, wt, out, info.State == `MERGED`); err != nil {
		return ``, err
	}
	if inside {
		return p.Root, nil
	}
	return dir, nil
}

// prPrune removes the worktrees of merged or closed PRs after asking.
func prPrune(dir string, in io.Reader, out io.Writer) error {
	p, err := openProject(dir)
	if err != nil {
		return err
	}
	wts, err := listWorktrees(p.Repo)
	if err != nil {
		return err
	}
	var found int
	var stale []worktree
	merged := map[string]bool{}
	for _, wt := range wts {
		n := prFromPath(p.Root, wt.Path)
		if n == `` {
			continue
		}
		found++
		state := `UNKNOWN`
		if info, err := ghPRView(p.Repo, n); err == nil {
			state = info.State
		}
		branch := wt.Branch
		if branch == `` {
			branch = `detached`
		}
		fmt.Fprintf(out, "  pr/%-6s  %-8s  %s\n", n, state, branch)
		if state == `MERGED` || state == `CLOSED` {
			stale = append(stale, wt)
			merged[wt.Path] = state == `MERGED`
		}
	}
	switch {
	case found == 0:
		fmt.Fprintln(out, `no PR worktrees found`)
		return nil
	case len(stale) == 0:
		fmt.Fprintln(out, `nothing to prune, all PR worktrees are open`)
		return nil
	}
	q := fmt.Sprintf(`%d worktree(s) can be pruned. Remove?`, len(stale))
	if !confirm(in, out, q) {
		fmt.Fprintln(out, `aborted`)
		return nil
	}
	var errs []error
	for _, wt := range stale {
		errs = append(errs, removeWorktree(p.Common, wt, out, merged[wt.Path]))
	}
	return errors.Join(errs...)
}

// removeWorktree removes a clean worktree and its local branch. It
// refuses a worktree with uncommitted or untracked files, and, unless
// the PR was merged, a branch with commits that exist nowhere else.
// Merged is exempt because squash merges leave such commits behind.
func removeWorktree(repo string, wt worktree, out io.Writer, merged bool) error {
	if wt.Branch != `` && !merged {
		n, err := git(repo, `rev-list`, `--count`, wt.Branch, `--not`,
			// With --branches, --exclude patterns omit the refs/heads/ prefix.
			`--exclude=`+wt.Branch, `--branches`, `--remotes`)
		if err != nil {
			return err
		}
		if n != `0` {
			return fmt.Errorf(`%s has %s unpushed commit(s); push them or delete the branch yourself`, wt.Branch, n)
		}
	}
	if _, err := git(repo, `worktree`, `remove`, wt.Path); err != nil {
		return fmt.Errorf(`%w (commit or stash first, or: git worktree remove --force %s)`, err, wt.Path)
	}
	fmt.Fprintf(out, "removed worktree %s\n", wt.Path)
	if wt.Branch == `` {
		return nil
	}
	if _, err := git(repo, `branch`, `-D`, wt.Branch); err != nil {
		return err
	}
	fmt.Fprintf(out, "deleted branch %s\n", wt.Branch)
	return nil
}

func findWorktree(wts []worktree, path string) (worktree, bool) {
	for _, wt := range wts {
		if samePath(wt.Path, path) {
			return wt, true
		}
	}
	return worktree{}, false
}

// prFromPath returns n when path is root/pr/<n>, else "".
func prFromPath(root, path string) string {
	rel, err := filepath.Rel(resolve(root), resolve(path))
	if err != nil {
		return ``
	}
	if m := prDir.FindStringSubmatch(filepath.ToSlash(rel)); m != nil {
		return m[1]
	}
	return ``
}

// prFromDir returns n when dir is inside the worktree root/pr/<n>.
func prFromDir(root, dir string) string {
	top, err := git(dir, `rev-parse`, `--show-toplevel`)
	if err != nil {
		return ``
	}
	return prFromPath(root, top)
}

func isInside(dir, root string) bool {
	d, r := resolve(dir), resolve(root)
	return d == r || strings.HasPrefix(d, r+string(filepath.Separator))
}

func confirm(in io.Reader, out io.Writer, question string) bool {
	fmt.Fprintf(out, "%s [y/N] ", question)
	line, _ := bufio.NewReader(in).ReadString('\n')
	return strings.EqualFold(strings.TrimSpace(line), `y`)
}
