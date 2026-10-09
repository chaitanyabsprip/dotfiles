// Package prompt provides small helpers for interactive CLI prompts.
package prompt

import (
	"bufio"
	"fmt"
	"io"
	"math/rand/v2"
	"strings"
)

// Confirm asks question on out and reads a y/N answer from in, returning
// true only for a "y" (case-insensitive) response. in is wrapped in a
// *bufio.Reader internally when it isn't already one; pass the same
// *bufio.Reader in and to a following ConfirmWord so no buffered input
// is lost between the two reads.
func Confirm(in io.Reader, out io.Writer, question string) bool {
	fmt.Fprintf(out, "%s [y/N] ", question)
	line, _ := reader(in).ReadString('\n')
	return strings.EqualFold(strings.TrimSpace(line), `y`)
}

const wordLetters = `abcdefghijklmnopqrstuvwxyz`

// ConfirmWord prints a random word on out and asks the user to type it
// back on in, returning true only on an exact match. Use as a harder,
// second confirmation before an irreversible action, passing the same
// *bufio.Reader given to a preceding Confirm.
func ConfirmWord(in io.Reader, out io.Writer) bool {
	word := randomWord(6)
	fmt.Fprintf(out, "Type %q to confirm: ", word)
	line, _ := reader(in).ReadString('\n')
	return strings.TrimSpace(line) == word
}

// reader returns in as a *bufio.Reader, wrapping it only if it isn't
// already one, so repeated calls sharing the same *bufio.Reader don't
// lose input buffered by an earlier call.
func reader(in io.Reader) *bufio.Reader {
	if r, ok := in.(*bufio.Reader); ok {
		return r
	}
	return bufio.NewReader(in)
}

func randomWord(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = wordLetters[rand.IntN(len(wordLetters))]
	}
	return string(b)
}
