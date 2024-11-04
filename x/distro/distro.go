package distro

import (
	"io"
	"log"
	"strings"

	"github.com/Chaitanyabsprip/dot/x/have"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai/run"
)

func Name() string {
	log.SetOutput(io.Discard)
	osReleaseFile := "/etc/os-release"
	if futil.Exists(osReleaseFile) {
		line, err := futil.FindString(osReleaseFile, "NAME=.*")
		if err != nil {
			return ""
		}
		return strings.Trim(
			strings.Split(line, "=")[1],
			`"`,
		)

	} else if ok, _ := have.Executable("lsb_release"); ok {
		return strings.TrimSpace(run.Out("lsb_release -si"))
	} else {
		return strings.TrimSpace(run.Out("uname -s"))
	}
}
