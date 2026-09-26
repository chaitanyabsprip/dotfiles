package workdirs

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseWorktrees(t *testing.T) {
	in := "worktree /p/.bare\nbare\n\n" +
		"worktree /p/main\nHEAD abc\nbranch refs/heads/main\n\n" +
		"worktree /p/pr/7\nHEAD def\ndetached\n"
	want := []worktree{
		{Path: `/p/.bare`, Bare: true},
		{Path: `/p/main`, Branch: `main`},
		{Path: `/p/pr/7`},
	}
	if got := parseWorktrees(in); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v\nwant %+v", got, want)
	}
}

func TestOpenProject(t *testing.T) {
	plain := newRepo(t, filepath.Join(t.TempDir(), `plain`))

	home := homeBase(t)
	gitT(t, filepath.Join(home, `root`), `worktree`, `add`, `-q`, `-b`, `feat/x`,
		filepath.Join(home, `feat`, `x`))

	bare := filepath.Join(t.TempDir(), `bare`)
	gitT(t, plain, `clone`, `-q`, `--bare`, plain, filepath.Join(bare, `.bare`))
	gitT(t, filepath.Join(bare, `.bare`), `worktree`, `add`, `-q`,
		filepath.Join(bare, `main`), `main`)

	tests := []struct{ dir, want string }{
		{plain, plain},
		{home, home}, // project dir itself is not a checkout
		{filepath.Join(home, `root`), home},
		{filepath.Join(home, `feat`, `x`), home},
		{bare, bare},
		{filepath.Join(bare, `.bare`), bare},
		{filepath.Join(bare, `main`), bare},
	}
	for _, tt := range tests {
		p, err := openProject(tt.dir)
		if err != nil {
			t.Errorf("openProject(%s): %v", tt.dir, err)
			continue
		}
		if !samePath(p.Root, tt.want) {
			t.Errorf("openProject(%s).Root = %s, want %s", tt.dir, p.Root, tt.want)
		}
		if _, err := git(p.Repo, `rev-parse`, `--git-dir`); err != nil {
			t.Errorf("openProject(%s).Repo %s is not usable: %v", tt.dir, p.Repo, err)
		}
	}

	if _, err := openProject(t.TempDir()); err == nil {
		t.Error("openProject outside any repo should fail")
	}
}
