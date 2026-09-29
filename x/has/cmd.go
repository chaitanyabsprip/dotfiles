package has

import (
	"fmt"
	"os"

	"github.com/rwxrob/bonzai"
)

var Cmd = &bonzai.Cmd{
	Name:    `has`,
	Short:   `find the nearest ancestor dir containing an entry`,
	Usage:   `[entry] [type]`,
	MaxArgs: 2,
	Long: `
Walks up from the current directory through its parents looking for
` + "`[entry]`" + ` (a file, directory, ...; default ` + "`.git`" + `, so bare
` + "`x has`" + ` finds the repository root), the way POSIX ` + "`test(1)`" + ` checks
it, and prints the path to the first match. ` + "`[type]`" + ` is a single
test-type character: ` + "`e`" + ` exists (default), ` + "`f`" + ` regular file,
` + "`d`" + ` directory, ` + "`L`" + ` symlink, or ` + "`s`" + ` non-empty file. Exits
non-zero with nothing printed when no ancestor matches.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		entry, typ := `.git`, byte('e')
		if len(args) >= 1 {
			entry = args[0]
		}
		if len(args) == 2 {
			if len(args[1]) != 1 {
				return fmt.Errorf(`invalid type %q, want a single character`, args[1])
			}
			typ = args[1][0]
		}
		path, err := Find(wd, entry, typ)
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	},
}
