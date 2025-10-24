package caseconv

import (
	"fmt"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
)

var Cmd = &bonzai.Cmd{
	Name:  "case",
	Short: "convert text to various cases",
	Opts: `lower|upper|camel|title|constant|header|sentence|snake|kebab`,
	Usage: `case [type] [text]`,
	Comp: comp.Opts,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		if len(args) < 2 {
			if len(args) < 1 {
				return fmt.Errorf("missing output type")
			}
			output := args[0]
			var b strings.Builder
			buf := make([]byte, 4096)
			for {
				n, err := fmt.Scanln(&buf)
				if n > 0 {
					b.Write(buf[:n])
					b.WriteByte('\n')
				}
				if err != nil {
					break
				}
			}
			text := b.String()
			if len(strings.TrimSpace(text)) == 0 {
				fmt.Println("No input provided via stdin.")
				return nil
			}
			args = []string{output, text}
		}
		output, text := args[0], strings.Join(args[1:], " ")
		text = strings.TrimSuffix(text, "\n")
		var result string
		switch output {
		case "lower":
			result = ToLower(text)
		case "upper":
			result = ToUpper(text)
		case "camel":
			result = ToCamel(text)
		case "title":
			result = ToTitle(text)
		case "constant":
			result = ToConstant(text)
		case "header":
			result = ToHeader(text)
		case "sentence":
			result = ToSentence(text)
		case "snake":
			result = ToSnake(text)
		case "kebab":
			result = ToKebab(text)
		default:
			fmt.Println("Unknown type:", output)
			return nil
		}
		fmt.Println(result)
		return nil
	},
}
