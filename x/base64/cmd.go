package base64

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
)

func isURLSafe() bool {
	val := os.Getenv("BASE64_URLSAFE")
	return val == "1" || strings.ToLower(val) == "true"
}

var EncodeCmd = &bonzai.Cmd{
	Name:  "encode",
	Short: "encode input to base64",
	Long: `
Encodes args, or stdin when no args are given, as base64. Set
BASE64_URLSAFE=1 to use the URL-safe alphabet instead of the standard
one.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		var input string
		if len(args) > 0 {
			input = strings.Join(args, " ")
		} else {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			input = string(data)
		}
		result, err := Encode(input, isURLSafe())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		fmt.Println(result)
		return nil
	},
}

var DecodeCmd = &bonzai.Cmd{
	Name:  "decode",
	Short: "decode base64 input",
	Long: `
Decodes args, or stdin when no args are given, from base64. Set
BASE64_URLSAFE=1 if the input used the URL-safe alphabet.`,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		var input string
		if len(args) > 0 {
			input = strings.Join(args, " ")
		} else {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			input = string(data)
		}
		result, err := Decode(strings.TrimSpace(input), isURLSafe())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		fmt.Print(result)
		return nil
	},
}

var Cmd = &bonzai.Cmd{
	Name:  "pem",
	Alias: `base64|b64`,
	Short: "base64 encode/decode utility",
	Long: `
Encodes or decodes text as base64. See 'pem help' for the two
commands.`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{EncodeCmd, DecodeCmd, help.Cmd},
	Do: func(x *bonzai.Cmd, _ ...string) error {
		fmt.Printf("%s - %s\n\n", x.Name, x.Short)
		fmt.Println(`COMMANDS:`)
		for _, c := range x.Cmds {
			fmt.Printf("  %-10s - %s\n", c.Name, c.Short)
		}
		return nil
	},
}
