package work

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/rwxrob/bonzai"
)

var approveCmd = &bonzai.Cmd{
	Name:    `approve`,
	Short:   `approve a github pull request`,
	Usage:   `<pr-url>`,
	NumArgs: 1,
	Long: `
Approves the pull request at ` + "`<pr-url>`" + `, a full GitHub PR URL such as
https://github.com/owner/repo/pull/123, via 'gh pr review --approve'.
Needs the gh CLI, authenticated with review rights on the repo.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		return approve(args[0])
	},
}

// prURL matches a GitHub pull request URL: github.com/<owner>/<repo>/pull/<number>.
var prURL = regexp.MustCompile(`^https://github\.com/[^/]+/[^/]+/pull/[0-9]+/?$`)

// approve validates url and, if it looks like a GitHub PR URL, approves it.
func approve(url string) error {
	if !prURL.MatchString(url) {
		return fmt.Errorf(`not a GitHub PR URL: %q`, url)
	}
	return ghApprove(url)
}

// ghApprove is a seam over the gh CLI; tests replace it.
var ghApprove = func(url string) error {
	cmd := exec.Command(`gh`, `pr`, `review`, url, `--approve`)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
