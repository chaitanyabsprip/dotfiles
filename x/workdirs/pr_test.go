package workdirs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeGH replaces the gh seams for one test and counts checkouts.
func fakeGH(t *testing.T, prs map[string]prInfo, checkout func(dir, n string) error) *int {
	t.Helper()
	view, co := ghPRView, ghPRCheckout
	t.Cleanup(func() { ghPRView, ghPRCheckout = view, co })
	calls := 0
	ghPRView = func(_, n string) (prInfo, error) {
		info, ok := prs[n]
		if !ok {
			return prInfo{}, errors.New(`no such PR`)
		}
		return info, nil
	}
	ghPRCheckout = func(dir, n string) error {
		calls++
		return checkout(dir, n)
	}
	return &calls
}

// createBranch emulates gh pr checkout by creating branch in the worktree.
func createBranch(branch string) func(dir, n string) error {
	return func(dir, _ string) error {
		_, err := git(dir, `switch`, `-q`, `-c`, branch)
		return err
	}
}

var noInput = strings.NewReader(``)

func TestPRCheckout(t *testing.T) {
	project := homeBase(t)
	calls := fakeGH(t, map[string]prInfo{`7`: {State: `OPEN`, HeadRefName: `feat/seven`}},
		createBranch(`feat/seven`))

	path, err := prCheckout(project, `7`, ``, noInput, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(project, `pr`, `7`); !samePath(path, want) {
		t.Errorf("path = %s, want %s", path, want)
	}
	if got := gitT(t, path, `branch`, `--show-current`); got != `feat/seven` {
		t.Errorf("branch = %q", got)
	}

	again, err := prCheckout(filepath.Join(project, `root`), `7`, ``, noInput, io.Discard)
	if err != nil || !samePath(again, path) {
		t.Errorf("second checkout = %s, %v; want %s", again, err, path)
	}
	if *calls != 1 {
		t.Errorf("gh pr checkout ran %d times, want 1", *calls)
	}
}

func TestPRCheckoutReusesBranchWorktree(t *testing.T) {
	project := homeBase(t)
	existing := filepath.Join(project, `feat`, `seven`)
	gitT(t, filepath.Join(project, `root`), `worktree`, `add`, `-q`, `-b`, `feat/seven`, existing)
	fakeGH(t, map[string]prInfo{`7`: {State: `OPEN`, HeadRefName: `feat/seven`}},
		func(string, string) error { return errors.New(`must not check out`) })

	path, err := prCheckout(project, `7`, ``, noInput, io.Discard)
	if err != nil || !samePath(path, existing) {
		t.Errorf("got %s, %v; want existing %s", path, err, existing)
	}
}

func TestPRCheckoutFailureCleansUp(t *testing.T) {
	project := homeBase(t)
	fakeGH(t, map[string]prInfo{`7`: {State: `OPEN`}},
		func(string, string) error { return errors.New(`boom`) })

	if _, err := prCheckout(project, `7`, ``, noInput, io.Discard); err == nil {
		t.Fatal(`want error`)
	}
	if _, err := os.Stat(filepath.Join(project, `pr`, `7`)); !os.IsNotExist(err) {
		t.Errorf("worktree left behind: %v", err)
	}
}

func TestPRCheckoutMergedNeedsConfirmation(t *testing.T) {
	project := homeBase(t)
	fakeGH(t, map[string]prInfo{`7`: {State: `MERGED`, HeadRefName: `feat/seven`}},
		createBranch(`feat/seven`))

	_, err := prCheckout(project, `7`, ``, strings.NewReader("n\n"), io.Discard)
	if !errors.Is(err, errAborted) {
		t.Errorf("answer n: err = %v, want errAborted", err)
	}
	if _, err := prCheckout(project, `7`, ``, strings.NewReader("y\n"), io.Discard); err != nil {
		t.Errorf("answer y: %v", err)
	}
}

func TestPRCheckoutRejectsBadNumber(t *testing.T) {
	project := homeBase(t)
	fakeGH(t, nil, createBranch(`x`))
	if _, err := prCheckout(project, `7; rm -rf ~`, ``, noInput, io.Discard); err == nil {
		t.Error(`non-numeric PR should be rejected`)
	}
}

func TestPRDone(t *testing.T) {
	project := homeBase(t)
	fakeGH(t, map[string]prInfo{`7`: {State: `OPEN`, HeadRefName: `feat/seven`}},
		createBranch(`feat/seven`))
	path, err := prCheckout(project, `7`, ``, noInput, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	cd, err := prDone(path, ``, ``, io.Discard) // PR inferred from inside pr/7
	if err != nil {
		t.Fatal(err)
	}
	if !samePath(cd, project) {
		t.Errorf("cd target = %s, want %s", cd, project)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("worktree still exists: %v", err)
	}
	if out := gitT(t, filepath.Join(project, `root`), `branch`, `--list`, `feat/seven`); out != `` {
		t.Errorf("branch not deleted: %q", out)
	}
}

// The shell passes $PWD unresolved; git records resolved paths
// (macOS /var -> /private/var). The cd target must still be the root.
func TestPRDoneFromUnresolvedPath(t *testing.T) {
	project := homeBase(t)
	fakeGH(t, map[string]prInfo{`7`: {State: `OPEN`, HeadRefName: `feat/seven`}},
		createBranch(`feat/seven`))
	if _, err := prCheckout(project, `7`, ``, noInput, io.Discard); err != nil {
		t.Fatal(err)
	}
	cd, err := prDone(filepath.Join(project, `pr`, `7`), ``, ``, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	if !samePath(cd, project) {
		t.Errorf("cd target = %s, want %s", cd, project)
	}
}

func TestPRDoneKeepsUnpushedCommits(t *testing.T) {
	project := homeBase(t)
	fakeGH(t, map[string]prInfo{`7`: {State: `OPEN`, HeadRefName: `feat/seven`}},
		createBranch(`feat/seven`))
	path, err := prCheckout(project, `7`, ``, noInput, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	gitT(t, path, `commit`, `-q`, `--allow-empty`, `-m`, `local only`)

	if _, err := prDone(project, `7`, ``, io.Discard); err == nil {
		t.Error(`done with unpushed commits should fail`)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("worktree removed: %v", err)
	}
	if out := gitT(t, filepath.Join(project, `root`), `branch`, `--list`, `feat/seven`); out == `` {
		t.Error(`branch with unpushed commits was deleted`)
	}
}

// gh checks out a fork's PR as a plain local branch with no
// remote-tracking ref. Its commits are on GitHub as the PR head, so done
// must not call them unpushed — but a local commit on top still is.
func TestPRDoneForkPR(t *testing.T) {
	project := homeBase(t)
	root := filepath.Join(project, `root`)
	fakeGH(t, nil, func(dir, _ string) error {
		if _, err := git(dir, `switch`, `-q`, `-c`, `feat/fork`); err != nil {
			return err
		}
		_, err := git(dir, `commit`, `-q`, `--allow-empty`, `-m`, `fork work`)
		return err
	})
	head := ``
	ghPRView = func(_, _ string) (prInfo, error) {
		return prInfo{State: `OPEN`, HeadRefName: `feat/fork`, HeadRefOid: head}, nil
	}

	path, err := prCheckout(project, `8`, ``, noInput, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	head = gitT(t, root, `rev-parse`, `feat/fork`) // what GitHub has

	gitT(t, path, `commit`, `-q`, `--allow-empty`, `-m`, `local only`)
	if _, err := prDone(project, `8`, ``, io.Discard); err == nil {
		t.Error(`a commit beyond the PR head should block done`)
	}

	gitT(t, path, `reset`, `-q`, `--hard`, head)
	if _, err := prDone(project, `8`, ``, io.Discard); err != nil {
		t.Errorf("clean fork PR worktree should be removable: %v", err)
	}
}

func TestPRDoneRejectsBadNumber(t *testing.T) {
	project := homeBase(t)
	other := filepath.Join(project, `feat`, `x`)
	gitT(t, filepath.Join(project, `root`), `worktree`, `add`, `-q`, `-b`, `feat/x`, other)
	fakeGH(t, nil, createBranch(`x`))

	if _, err := prDone(project, `../feat/x`, ``, io.Discard); err == nil {
		t.Error(`non-numeric PR should be rejected`)
	}
	if _, err := os.Stat(other); err != nil {
		t.Errorf("unrelated worktree removed: %v", err)
	}
}

func TestPRDoneKeepsDirtyWorktree(t *testing.T) {
	project := homeBase(t)
	fakeGH(t, map[string]prInfo{`7`: {State: `OPEN`, HeadRefName: `feat/seven`}},
		createBranch(`feat/seven`))
	path, err := prCheckout(project, `7`, ``, noInput, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	wip := filepath.Join(path, `wip.txt`)
	if err := os.WriteFile(wip, []byte(`unsaved`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := prDone(project, `7`, ``, io.Discard); err == nil {
		t.Error(`done on a dirty worktree should fail`)
	}
	if _, err := os.Stat(wip); err != nil {
		t.Errorf("uncommitted work was deleted: %v", err)
	}
}

func TestPRPrune(t *testing.T) {
	project := homeBase(t)
	fakeGH(t, map[string]prInfo{
		`7`: {State: `OPEN`, HeadRefName: `b7`},
		`8`: {State: `OPEN`, HeadRefName: `b8`},
	}, func(dir, n string) error { _, err := git(dir, `switch`, `-q`, `-c`, `b`+n); return err })
	p7, err := prCheckout(project, `7`, ``, noInput, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	p8, err := prCheckout(project, `8`, ``, noInput, io.Discard)
	if err != nil {
		t.Fatal(err)
	}

	// Squash-merged PRs leave local-only commits; merged still prunes.
	gitT(t, p7, `commit`, `-q`, `--allow-empty`, `-m`, `squashed upstream`)
	ghPRView = func(_, n string) (prInfo, error) {
		if n == `7` {
			return prInfo{State: `MERGED`}, nil
		}
		return prInfo{State: `OPEN`}, nil
	}
	if err := prPrune(project, strings.NewReader("y\n"), io.Discard); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p7); !os.IsNotExist(err) {
		t.Errorf("merged PR worktree not pruned: %v", err)
	}
	if _, err := os.Stat(p8); err != nil {
		t.Errorf("open PR worktree was pruned: %v", err)
	}
}
