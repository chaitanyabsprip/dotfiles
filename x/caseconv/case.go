// Package caseconv provides command to convert strings to various cases.
package caseconv

import (
	"strings"
	"unicode"
)

func splitWords(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToCamel(s string) string {
	words := splitWords(s)
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(strings.ToLower(words[0]))
	for _, w := range words[1:] {
		if len(w) > 0 {
			b.WriteString(strings.ToUpper(string(w[0])))
			b.WriteString(strings.ToLower(w[1:]))
		}
	}
	return b.String()
}

func ToTitle(s string) string {
	words := splitWords(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(
				string(w[0]),
			) + strings.ToLower(
				w[1:],
			)
		}
	}
	return strings.Join(words, " ")
}

func ToConstant(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToUpper(w)
	}
	return strings.Join(words, "_")
}

func ToHeader(s string) string {
	words := splitWords(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(
				string(w[0]),
			) + strings.ToLower(
				w[1:],
			)
		}
	}
	return strings.Join(words, "-")
}

func ToSentence(s string) string {
	s = strings.ToLower(s)
	for i, r := range s {
		if unicode.IsLetter(r) {
			return string(unicode.ToUpper(r)) + s[i+1:]
		}
	}
	return s
}

func ToSnake(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}

func ToKebab(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "-")
}

