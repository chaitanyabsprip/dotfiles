package work

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/rwxrob/bonzai"
)

var approveCmd = &bonzai.Cmd{
	Name:    `approve`,
	Short:   `approve a github pull request`,
	Usage:   `<pr-url>|<pr-number>`,
	NumArgs: 1,
	Long: `
Approves the pull request at ` + "`<pr-url>`" + `, a full GitHub PR URL such as
https://github.com/owner/repo/pull/123, via 'gh pr review --approve'.
A bare ` + "`<pr-number>`" + ` is taken as a PR in the GitHub repo of the
current directory, and fails outside one.
Needs the gh CLI, authenticated with review rights on the repo.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		return approve(args[0])
	},
}

// prURL matches a GitHub pull request URL: github.com/<owner>/<repo>/pull/<number>.
var prURL = regexp.MustCompile(`^https://github\.com/[^/]+/[^/]+/pull/[0-9]+/?$`)

// approve validates url and, if it looks like a GitHub PR URL, approves it.
// A bare PR number is expanded against the current directory's GitHub repo.
func approve(url string) error {
	if prNumber.MatchString(url) {
		repo, err := ghRepoURL()
		if err != nil {
			return fmt.Errorf(`PR number %s needs a GitHub repo in the current directory: %w`, url, err)
		}
		url = repo + `/pull/` + url
	}
	if !prURL.MatchString(url) {
		return fmt.Errorf(`not a GitHub PR URL: %q`, url)
	}
	return ghApprove(url)
}

// Seams over the gh CLI; tests replace them.
var (
	ghApprove = func(url string) error {
		cmd := exec.Command(`gh`, `pr`, `review`, url, `--approve`)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		return cmd.Run()
	}
	// ghRepoURL returns the current directory's repo URL, e.g. https://github.com/owner/repo.
	ghRepoURL = func() (string, error) {
		out, err := exec.Command(`gh`, `repo`, `view`, `--json`, `url`, `-q`, `.url`).Output()
		return strings.TrimSpace(string(out)), err
	}
)
