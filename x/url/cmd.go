package url

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `url`,
	Short: `percent-encode or decode text`,
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{encodeCmd, decodeCmd},
}

var encodeCmd = &bonzai.Cmd{
	Name:  `encode`,
	Alias: `e`,
	Short: `percent-encode args or each stdin line`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		return convert(args, os.Stdin, os.Stdout, func(s string) (string, error) {
			return Encode(s), nil
		})
	},
}

var decodeCmd = &bonzai.Cmd{
	Name:  `decode`,
	Alias: `d`,
	Short: `decode args or each stdin line`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		return convert(args, os.Stdin, os.Stdout, Decode)
	},
}

// convert applies fn to the joined args, or to each line of r when there
// are no args. Lines are converted one at a time so a newline is never
// encoded into the output.
func convert(
	args []string,
	r io.Reader,
	w io.Writer,
	fn func(string) (string, error),
) error {
	if len(args) > 0 {
		out, err := fn(strings.Join(args, ` `))
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, out)
		return err
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		out, err := fn(sc.Text())
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, out); err != nil {
			return err
		}
	}
	return sc.Err()
}
