package json

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chaitanyabsprip/bonzai"
	"github.com/chaitanyabsprip/bonzai/cmds/help"
	"github.com/chaitanyabsprip/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  `json`,
	Short: `json helpers`,
	Long: `
Small JSON helpers. See 'json help' for the commands.`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{unescapeCmd, help.Cmd},
	Do: func(x *bonzai.Cmd, _ ...string) error {
		fmt.Printf("%s - %s\n\n", x.Name, x.Short)
		fmt.Println(`COMMANDS:`)
		for _, c := range x.Cmds {
			fmt.Printf("  %-10s - %s\n", c.Name, c.Short)
		}
		return nil
	},
}

var unescapeCmd = &bonzai.Cmd{
	Name:  `unescape`,
	Alias: `u`,
	Short: `expand JSON strings nested inside JSON`,
	Long: `
Reads JSON from args, or stdin when no args are given, and replaces
every string holding a JSON object or array with the parsed value,
recursively. An envelope carrying a stringified payload comes out as
one nested, indented document.

  tug peek my-topic | x json unescape`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		var input []byte
		if len(args) > 0 {
			input = []byte(strings.Join(args, ` `))
		} else {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			input = data
		}
		out, err := Unescape(input)
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}
