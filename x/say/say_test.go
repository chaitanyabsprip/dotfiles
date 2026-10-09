package say

import "testing"

// Color escapes collapse to "" outside a terminal (term.IsInteractive),
// which `go test` always is, so these assert the plain text: symbol,
// spacing and message, matching the original shell scripts' printf.
func TestStyles(t *testing.T) {
	tests := []struct {
		name string
		fn   func(string) string
		want string
	}{
		{`Error`, Error, " ❌ boom"},
		{`Warning`, Warning, " ⚠️ boom"},
		{`Success`, Success, " ✔ boom"},
		{`InProgress`, InProgress, " ... boom"},
		{`Bold`, Bold, `boom`},
		{`Italic`, Italic, `boom`},
	}
	for _, tt := range tests {
		if got := tt.fn(`boom`); got != tt.want {
			t.Errorf("%s(%q) = %q, want %q", tt.name, `boom`, got, tt.want)
		}
	}
}
