package workdirs

import (
	"path/filepath"
	"testing"
)

func TestWorktreePathFor(t *testing.T) {
	tests := map[string]string{
		`feat/mod/TICKET-1/short-desc`: `feat/short-desc`,
		`fix/typo`:                     `fix/typo`,
		`hotfix`:                       `hotfix`,
		`/feat/x/`:                     `feat/x`,
	}
	for in, want := range tests {
		if got := worktreePathFor(in); got != filepath.FromSlash(want) {
			t.Errorf("worktreePathFor(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAddWorktree(t *testing.T) {
	project := homeBase(t) // unresolved temp path on purpose (Review Focus)
	root := filepath.Join(project, `root`)

	path, err := addWorktree(root, `feat/mod/T-1/new-thing`)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(project, `feat`, `new-thing`); !samePath(path, want) {
		t.Errorf("path = %s, want %s", path, want)
	}
	if got := gitT(t, path, `branch`, `--show-current`); got != `feat/mod/T-1/new-thing` {
		t.Errorf("branch = %q", got)
	}

	again, err := addWorktree(project, `feat/mod/T-1/new-thing`)
	if err != nil || !samePath(again, path) {
		t.Errorf("second add = %s, %v; want existing %s", again, err, path)
	}

	gitT(t, root, `branch`, `fix/typo`)
	path, err = addWorktree(root, `fix/typo`)
	if err != nil {
		t.Fatal(err)
	}
	if got := gitT(t, path, `branch`, `--show-current`); got != `fix/typo` {
		t.Errorf("existing branch not checked out: %q", got)
	}
}

func TestAddWorktreeRemoteOnlyBranch(t *testing.T) {
	upstream := newRepo(t, filepath.Join(t.TempDir(), `upstream`))
	gitT(t, upstream, `branch`, `feat/remote`)
	project := filepath.Join(t.TempDir(), `project`)
	gitT(t, upstream, `clone`, `-q`, upstream, filepath.Join(project, `root`))

	path, err := addWorktree(filepath.Join(project, `root`), `feat/remote`)
	if err != nil {
		t.Fatal(err)
	}
	got := gitT(t, path, `rev-parse`, `--abbrev-ref`, `@{upstream}`)
	if got != `origin/feat/remote` {
		t.Errorf("upstream = %q, want origin/feat/remote", got)
	}
}
