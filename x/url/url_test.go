package url

import (
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct{ in, want string }{
		{`hello world`, `hello%20world`},
		{`a+b=c&d`, `a%2Bb%3Dc%26d`},
		{`/path?q#f`, `%2Fpath%3Fq%23f`},
		{`-_.~AZaz09`, `-_.~AZaz09`},
		{`é`, `%C3%A9`},
		{``, ``},
	}
	for _, tt := range tests {
		if got := Encode(tt.in); got != tt.want {
			t.Errorf("Encode(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct{ in, want string }{
		{`hello%20world`, `hello world`},
		{`a+b`, `a b`},
		{`%C3%A9`, `é`},
		{`%2fx`, `/x`},
	}
	for _, tt := range tests {
		got, err := Decode(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("Decode(%q) = %q, %v; want %q", tt.in, got, err, tt.want)
		}
	}
	if _, err := Decode(`100%zz`); err == nil {
		t.Error(`Decode("100%zz") should fail`)
	}
}

func FuzzRoundTrip(f *testing.F) {
	f.Add(`hello world`)
	f.Add(`a+b/c?d=e&f`)
	f.Add(`é😀`)
	f.Fuzz(func(t *testing.T, s string) {
		got, err := Decode(Encode(s))
		if err != nil || got != s {
			t.Errorf("Decode(Encode(%q)) = %q, %v", s, got, err)
		}
	})
}

func TestConvert(t *testing.T) {
	enc := func(s string) (string, error) { return Encode(s), nil }

	var out strings.Builder
	in := strings.NewReader("a b\n\nc/d") // blank line, no trailing newline
	if err := convert(nil, in, &out, enc); err != nil {
		t.Fatal(err)
	}
	if want := "a%20b\n\nc%2Fd\n"; out.String() != want {
		t.Errorf("stdin: got %q, want %q", out.String(), want)
	}

	out.Reset()
	if err := convert([]string{`a`, `b`}, strings.NewReader("ignored\n"), &out, enc); err != nil {
		t.Fatal(err)
	}
	if want := "a%20b\n"; out.String() != want {
		t.Errorf("args: got %q, want %q", out.String(), want)
	}

	if err := convert(nil, strings.NewReader("ok\n%zz\n"), &out, Decode); err == nil {
		t.Error("bad escape on stdin should fail")
	}
}
