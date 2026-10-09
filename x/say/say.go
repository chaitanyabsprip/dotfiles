// Package say renders the short colored status messages that other
// scripts in this repo print through: error, warning, success,
// in-progress, bold and italic. Colors collapse to plain text when
// stdout isn't a terminal (see bonzai/term).
package say

import "github.com/rwxrob/bonzai/term"

func Error(msg string) string      { return colorize(term.Bold+term.Red, " ❌ "+msg) }
func Warning(msg string) string    { return colorize(term.Bold+term.Yellow, " ⚠️ "+msg) }
func Success(msg string) string    { return colorize(term.Bold+term.Green, " ✔ "+msg) }
func InProgress(msg string) string { return colorize(term.Bold+term.Blue, " ... "+msg) }
func Bold(msg string) string       { return colorize(term.Bold, msg) }
func Italic(msg string) string     { return colorize(term.Italic, msg) }

func colorize(attr, s string) string { return attr + s + term.Reset }
