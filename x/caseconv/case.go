// Package caseconv provides command to convert strings to various cases.
package caseconv

import (
	"strings"
	"unicode"
)

func ToLower(s string) string {
	var b strings.Builder
	for _, r := range s {
		b.WriteRune(transform(r, unicode.ToLower))
	}
	return b.String()
}

func ToUpper(s string) string {
	var b strings.Builder
	for _, r := range s {
		b.WriteRune(transform(r, unicode.ToUpper))
	}
	return b.String()
}

func ToCamel(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	// Write first word in lower case, preserving non-letters
	for _, r := range words[0] {
		b.WriteRune(transform(r, unicode.ToLower))
	}
	for _, w := range words[1:] {
		if len(w) > 0 {
			for i, r := range w {
				var fn func(rune) rune
				if i == 0 {
					fn = unicode.ToUpper
				} else {
					fn = unicode.ToLower
				}
				b.WriteRune(transform(r, fn))
			}
		}
	}
	return b.String()
}

func ToTitle(s string) string {
	words := splitWords(s)
	for i, w := range words {
		if len(w) > 0 {
			var b strings.Builder
			for j, r := range w {
				var fn func(rune) rune
				if j == 0 {
					fn = unicode.ToUpper
				} else {
					fn = unicode.ToLower
				}
				b.WriteRune(transform(r, fn))
			}
			words[i] = b.String()
		}
	}
	return strings.Join(words, " ")
}

func ToConstant(s string) string {
	words := splitWords(s)
	for i, w := range words {
		var b strings.Builder
		for _, r := range w {
			b.WriteRune(transform(r, unicode.ToUpper))
		}
		words[i] = b.String()
	}
	return strings.Join(words, "_")
}

func ToHeader(s string) string {
	words := splitWords(s)
	for i, w := range words {
		if len(w) > 0 {
			var b strings.Builder
			for j, r := range w {
				var fn func(rune) rune
				if j == 0 {
					fn = unicode.ToUpper
				} else {
					fn = unicode.ToLower
				}
				b.WriteRune(transform(r, fn))
			}
			words[i] = b.String()
		}
	}
	return strings.Join(words, "-")
}

func ToSentence(s string) string {
	var b strings.Builder
	first := true
	for _, r := range s {
		if first && unicode.IsLetter(r) {
			b.WriteRune(unicode.ToUpper(r))
			first = false
		} else if unicode.IsLetter(r) {
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func ToSnake(s string) string {
	words := splitWords(s)
	for i, w := range words {
		var b strings.Builder
		for _, r := range w {
			b.WriteRune(transform(r, unicode.ToLower))
		}
		words[i] = b.String()
	}
	return strings.Join(words, "_")
}

func ToKebab(s string) string {
	words := splitWords(s)
	for i, w := range words {
		var b strings.Builder
		for _, r := range w {
			b.WriteRune(transform(r, unicode.ToLower))
		}
		words[i] = b.String()
	}
	return strings.Join(words, "-")
}

// splitWords splits on non-alphanumerics and on case boundaries:
// "fooBar" -> foo Bar, "HTTPServer" -> HTTP Server, "v2Api" -> v2 Api.
func splitWords(s string) []string {
	var words []string
	for _, field := range strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) {
		rs := []rune(field)
		start := 0
		for i := 1; i < len(rs); i++ {
			prev, cur := rs[i-1], rs[i]
			lowerToUpper := !unicode.IsUpper(prev) && unicode.IsUpper(cur)
			acronymEnd := unicode.IsUpper(prev) && unicode.IsUpper(cur) &&
				i+1 < len(rs) && unicode.IsLower(rs[i+1])
			if lowerToUpper || acronymEnd {
				words = append(words, string(rs[start:i]))
				start = i
			}
		}
		words = append(words, string(rs[start:]))
	}
	return words
}

func transform(r rune, fn func(rune) rune) rune {
	if unicode.IsLetter(r) {
		return fn(r)
	}
	return r
}
