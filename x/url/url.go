// Package url percent-encodes and decodes text using the RFC 3986
// unreserved set, so the output is safe anywhere in a URL.
package url

import (
	neturl "net/url"
	"strings"
)

// Encode escapes every byte except A-Z a-z 0-9 - _ . ~ and writes a
// space as %20 rather than +.
func Encode(s string) string {
	return strings.ReplaceAll(neturl.QueryEscape(s), `+`, `%20`)
}

// Decode reverses Encode. It also reads + as a space, as form data does.
func Decode(s string) (string, error) {
	return neturl.QueryUnescape(s)
}
