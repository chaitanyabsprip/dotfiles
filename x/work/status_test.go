package work

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStatusesListsEachWorktree(t *testing.T) {
	project := homeBase(t)
	root := filepath.Join(project, `root`)
	dirty := filepath.Join(project, `feat`, `dirty`)
	gitT(t, root, `worktree`, `add`, `-q`, `-b`, `feat/dirty`, dirty)

	if err := os.WriteFile(filepath.Join(dirty, `wip.txt`), []byte(`x`), 0o644); err != nil {
		t.Fatal(err)
	}

	sts, err := statuses(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(sts) != 2 {
		t.Fatalf("got %d worktrees, want 2: %+v", len(sts), sts)
	}

	byBranch := map[string]worktreeStatus{}
	for _, s := range sts {
		byBranch[s.Branch] = s
	}

	main, ok := byBranch[`main`]
	if !ok || !main.Status.IsClean {
		t.Errorf("main worktree not reported clean: %+v", main)
	}

	feat, ok := byBranch[`feat/dirty`]
	if !ok || feat.Status.IsClean || feat.Status.NumUntracked != 1 {
		t.Errorf("feat/dirty worktree not reported with 1 untracked file: %+v", feat)
	}
	if !samePath(feat.Path, dirty) {
		t.Errorf("feat/dirty path = %q, want %q", feat.Path, dirty)
	}
}

func TestStatusesIncludesPRState(t *testing.T) {
	project := homeBase(t)
	root := filepath.Join(project, `root`)
	pr := filepath.Join(project, `pr`, `8`)
	gitT(t, root, `worktree`, `add`, `-q`, `-b`, `feat/eight`, pr)
	fakeGH(t, map[string]prInfo{`8`: {State: `MERGED`}}, nil)

	sts, err := statuses(root)
	if err != nil {
		t.Fatal(err)
	}
	byBranch := map[string]worktreeStatus{}
	for _, s := range sts {
		byBranch[s.Branch] = s
	}

	pw, ok := byBranch[`feat/eight`]
	if !ok || pw.PR != `8` || pw.PRState != `MERGED` {
		t.Errorf("pr worktree = %+v, want PR=8 PRState=MERGED", pw)
	}

	main := byBranch[`main`]
	if main.PR != `` {
		t.Errorf("non-PR worktree got PR=%q, want \"\"", main.PR)
	}
}

func TestStatusesUnknownWhenGHFails(t *testing.T) {
	project := homeBase(t)
	root := filepath.Join(project, `root`)
	pr := filepath.Join(project, `pr`, `9`)
	gitT(t, root, `worktree`, `add`, `-q`, `-b`, `feat/nine`, pr)
	// No fakeGH: ghPRView is the TestMain fail-closed default, which errors.

	sts, err := statuses(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range sts {
		if s.Branch == `feat/nine` {
			if s.PR != `9` || s.PRState != `UNKNOWN` {
				t.Errorf("pr worktree with failing gh = %+v, want PR=9 PRState=UNKNOWN", s)
			}
			return
		}
	}
	t.Fatal(`feat/nine worktree not found`)
}

func TestStatusesSkipsBareEntry(t *testing.T) {
	src := newRepo(t, filepath.Join(t.TempDir(), `src`))
	project := filepath.Join(t.TempDir(), `project`)
	gitT(t, src, `clone`, `-q`, `--bare`, src, filepath.Join(project, `.bare`))
	main := filepath.Join(project, `main`)
	gitT(t, filepath.Join(project, `.bare`), `worktree`, `add`, `-q`, main, `main`)

	sts, err := statuses(main)
	if err != nil {
		t.Fatal(err)
	}
	if len(sts) != 1 || !samePath(sts[0].Path, main) {
		t.Errorf("statuses = %+v, want exactly the main worktree", sts)
	}
}

func TestStatusesOutsideProjectErrors(t *testing.T) {
	if _, err := statuses(t.TempDir()); err == nil {
		t.Error(`statuses outside a git project should error`)
	}
}
