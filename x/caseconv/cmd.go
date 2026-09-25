package caseconv

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
	Name:  "case",
	Short: "convert text to various cases",
	Opts:  `lower|upper|camel|title|constant|header|sentence|snake|kebab`,
	Usage: `case <type> [text]  (reads stdin line by line when text is omitted)`,
	Comp:  comp.Opts,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing output type")
		}
		if len(args) == 1 {
			return convertLines(args[0], os.Stdin, os.Stdout)
		}
		result, err := convert(args[0], strings.Join(args[1:], " "))
		if err != nil {
			return err
		}
		fmt.Println(result)
		return nil
	},
}

// convertLines converts each line of r independently, like a Unix filter.
func convertLines(kind string, r io.Reader, w io.Writer) error {
	if _, err := convert(kind, ""); err != nil {
		return err
	}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		result, _ := convert(kind, sc.Text())
		fmt.Fprintln(w, result)
	}
	return sc.Err()
}

func convert(kind, text string) (string, error) {
	switch kind {
	case "lower":
		return ToLower(text), nil
	case "upper":
		return ToUpper(text), nil
	case "camel":
		return ToCamel(text), nil
	case "title":
		return ToTitle(text), nil
	case "constant":
		return ToConstant(text), nil
	case "header":
		return ToHeader(text), nil
	case "sentence":
		return ToSentence(text), nil
	case "snake":
		return ToSnake(text), nil
	case "kebab":
		return ToKebab(text), nil
	}
	return "", fmt.Errorf("unknown type: %s", kind)
}
