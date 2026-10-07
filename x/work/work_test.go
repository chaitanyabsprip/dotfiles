package work

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

// TestWorktreesFindsSiblingOfBareRoot guards against the sessionizer
// bug where a bare-repo project's own ".git" pointer file (a sibling
// of its worktree directories, e.g. "ss/.git" next to "ss/main")
// caused the filesystem walk to skip the rest of that directory,
// silently hiding every real worktree next to it.
func TestWorktreesFindsSiblingOfBareRoot(t *testing.T) {
	isolateGit(t)
	tmp := t.TempDir()

	src := newRepo(t, filepath.Join(tmp, `src`))

	projectDir := filepath.Join(tmp, `projects`, `ss`)
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatal(err)
	}
	gitT(t, tmp, `clone`, `-q`, `--bare`, src, filepath.Join(projectDir, `.bare`))
	if err := os.WriteFile(
		filepath.Join(projectDir, `.git`), []byte("gitdir: ./.bare\n"), 0o644,
	); err != nil {
		t.Fatal(err)
	}
	gitT(t, projectDir, `worktree`, `add`, `-q`, `main`)

	orig := env.Projects
	env.Projects = filepath.Join(tmp, `projects`)
	t.Cleanup(func() { env.Projects = orig })

	got := Worktrees()
	want := filepath.Join(projectDir, `main`)
	for _, w := range got {
		if w == want {
			return
		}
	}
	t.Errorf("Worktrees() = %v, want to contain %q", got, want)
}
