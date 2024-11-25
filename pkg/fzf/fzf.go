package fzf

import (
	"bytes"
	"os/exec"
	"strings"
)

// func Select(in, opts []string) (out string, errO error) {
// 	inputChan := make(chan string)
// 	go func() {
// 		for _, s := range in {
// 			inputChan <- s
// 		}
// 		close(inputChan)
// 	}()
// 	outputChan := make(chan string)
// 	options, err := fzf.ParseOptions(true, append(opts, "--multi=1"))
// 	if err != nil {
// 		return "", fmt.Errorf(
// 			"fzf exited with code %d: %w",
// 			fzf.ExitError,
// 			err,
// 		)
// 	}
// 	options.Input = inputChan
// 	options.Output = outputChan
// 	fmt.Printf("%+v\n", options)
// 	code, err := fzf.Run(options)
// 	if err != nil {
// 		return "", fmt.Errorf(
// 			"fzf exited with code %d: %w",
// 			code,
// 			err,
// 		)
// 	}
// 	out = <-outputChan
// 	return
// }

func Select(in []string, fzfOpts ...string) (out string, errO error) {
	cmd := exec.Command("fzf", fzfOpts...)
	cmd.Stdin = strings.NewReader(strings.Join(in, "\n"))
	var outb bytes.Buffer
	cmd.Stdout = &outb
	err := cmd.Run()
	if err != nil {
		errO = err
	}
	out = strings.TrimSpace(outb.String())
	return
}

func SelectMulti(in ...string) []string {
	return []string{}
}
