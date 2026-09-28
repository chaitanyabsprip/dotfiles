// Package version holds the build-time version string shared by the dot
// and x binaries. Version is set via -ldflags at build time (see
// .github/workflows/build.yml and release.yml) to the tag exactly on the
// built commit, or that commit's short hash if no tag is on it — never
// git describe's "<tag>-<N>-g<hash>" form. Left at its zero value for a
// plain `go build`/`go run` with no ldflags.
package version

import (
	"fmt"

	"github.com/rwxrob/bonzai"
)

var Version = `dev`

// Cmd prints Version. Composed into both dot.Cmd (cmd.go) and x.Cmd
// (x/x.go) as `version`.
var Cmd = &bonzai.Cmd{
	Name:  `version`,
	Short: `print the build version`,
	Do: func(_ *bonzai.Cmd, _ ...string) error {
		fmt.Println(Version)
		return nil
	},
}
