package base64

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

func isURLSafe() bool {
	val := os.Getenv("BASE64_URLSAFE")
	return val == "1" || strings.ToLower(val) == "true"
}

var EncodeCmd = &bonzai.Cmd{
	Name:  "encode",
	Short: "encode input to base64",
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
	Comp:  comp.Cmds,
	Cmds:  []*bonzai.Cmd{EncodeCmd, DecodeCmd},
}
