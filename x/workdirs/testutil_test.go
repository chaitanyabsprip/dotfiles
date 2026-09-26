package workdirs

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// isolateGit keeps the user's git config (signing, hooks, templates) out
// of tests and gives commits a fixed identity.
func isolateGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath(`git`); err != nil {
		t.Skip(`git not installed`)
	}
	t.Setenv(`GIT_CONFIG_GLOBAL`, os.DevNull)
	t.Setenv(`GIT_CONFIG_NOSYSTEM`, `1`)
	for _, k := range []string{`GIT_AUTHOR_NAME`, `GIT_COMMITTER_NAME`} {
		t.Setenv(k, `test`)
	}
	for _, k := range []string{`GIT_AUTHOR_EMAIL`, `GIT_COMMITTER_EMAIL`} {
		t.Setenv(k, `test@example.com`)
	}
}

// gitT runs git in dir and fails the test on error.
func gitT(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := git(dir, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// newRepo creates a repository at dir with one commit on main.
func newRepo(t *testing.T, dir string) string {
	t.Helper()
	isolateGit(t)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	gitT(t, dir, `init`, `-q`, `-b`, `main`)
	gitT(t, dir, `commit`, `-q`, `--allow-empty`, `-m`, `init`)
	return dir
}

// homeBase creates the home-base layout and returns the project dir:
// project/root is the main checkout, worktrees go next to it.
func homeBase(t *testing.T) string {
	t.Helper()
	project := filepath.Join(t.TempDir(), `project`)
	newRepo(t, filepath.Join(project, `root`))
	return project
}
