package workdirs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepairProjectAfterMove(t *testing.T) {
	old := homeBase(t)
	gitT(t, filepath.Join(old, `root`), `worktree`, `add`, `-q`, `-b`, `feat/x`,
		filepath.Join(old, `feat`, `x`))

	// A submodule-style .git file must be left alone (Review Focus).
	sub := filepath.Join(old, `root`, `vendor`, `lib`)
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	gitfile := []byte("gitdir: ../../.git/modules/lib\n")
	if err := os.WriteFile(filepath.Join(sub, `.git`), gitfile, 0o644); err != nil {
		t.Fatal(err)
	}

	// A submodule inside a linked worktree points under worktrees/x/modules/;
	// it is not a worktree either.
	wtSub := filepath.Join(old, `feat`, `x`, `lib`)
	if err := os.MkdirAll(wtSub, 0o755); err != nil {
		t.Fatal(err)
	}
	wtGitfile := []byte("gitdir: ../../../root/.git/worktrees/x/modules/lib\n")
	if err := os.WriteFile(filepath.Join(wtSub, `.git`), wtGitfile, 0o644); err != nil {
		t.Fatal(err)
	}

	moved := old + `-moved`
	if err := os.Rename(old, moved); err != nil {
		t.Fatal(err)
	}
	wt := filepath.Join(moved, `feat`, `x`)
	if _, err := git(wt, `status`); err == nil {
		t.Fatal(`worktree should be broken after the move`)
	}

	linked, err := repairProject(moved)
	if err != nil {
		t.Fatal(err)
	}
	if len(linked) != 1 || !samePath(linked[0], wt) {
		t.Errorf("linked = %v, want [%s]", linked, wt)
	}
	if got := gitT(t, wt, `branch`, `--show-current`); got != `feat/x` {
		t.Errorf("branch after repair = %q", got)
	}
}

func TestRepairBareProjectAfterMove(t *testing.T) {
	src := newRepo(t, filepath.Join(t.TempDir(), `src`))
	old := filepath.Join(t.TempDir(), `project`)
	gitT(t, src, `clone`, `-q`, `--bare`, src, filepath.Join(old, `.bare`))
	gitT(t, filepath.Join(old, `.bare`), `worktree`, `add`, `-q`,
		filepath.Join(old, `main`), `main`)

	moved := old + `-moved`
	if err := os.Rename(old, moved); err != nil {
		t.Fatal(err)
	}
	if _, err := repairProject(moved); err != nil {
		t.Fatal(err)
	}
	gitT(t, filepath.Join(moved, `main`), `status`)
}

func TestRepairProjectNoRepo(t *testing.T) {
	if _, err := repairProject(t.TempDir()); err == nil {
		t.Error(`repair outside a project should fail`)
	}
}
