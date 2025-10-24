package caseconv

import (
	"testing"
	"unicode/utf8"
)

func FuzzToLower(f *testing.F) {
	f.Add("Hello World")
	f.Add("")
	f.Add("123-abc_!@#")
	f.Add("こんにちは世界")
	f.Add("😀😃😄")
	f.Fuzz(func(t *testing.T, input string) {
		out := ToLower(input)
		if !utf8.ValidString(out) {
			t.Errorf("output is not valid UTF-8: %q", out)
		}
	})
}

func FuzzToUpper(f *testing.F) {
	f.Add("Hello World")
	f.Add("")
	f.Add("123-abc_!@#")
	f.Add("こんにちは世界")
	f.Add("😀😃😄")
	f.Fuzz(func(t *testing.T, input string) {
		out := ToUpper(input)
		if !utf8.ValidString(out) {
			t.Errorf("output is not valid UTF-8: %q", out)
		}
	})
}

func FuzzToCamel(f *testing.F) {
	f.Add("hello world")
	f.Add("HELLO_WORLD")
	f.Add("123-abc")
	f.Add("")
	f.Add("こんにちは 世界")
	f.Add("😀😃😄")
	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic for input: %q", input)
			}
		}()
		out := ToCamel(input)
		if !utf8.ValidString(out) {
			t.Errorf("output is not valid UTF-8: %q", out)
		}
	})
}

func FuzzToTitle(f *testing.F) {
	f.Add("hello world")
	f.Add("HELLO_WORLD")
	f.Add("123-abc")
	f.Add("")
	f.Add("こんにちは 世界")
	f.Add("😀😃😄")
	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic for input: %q", input)
			}
		}()
		out := ToTitle(input)
		if !utf8.ValidString(out) {
			t.Errorf("output is not valid UTF-8: %q", out)
		}
	})
}

func FuzzToConstant(f *testing.F) {
	f.Add("hello world")
	f.Add("HELLO_WORLD")
	f.Add("123-abc")
	f.Add("")
	f.Add("こんにちは 世界")
	f.Add("😀😃😄")
	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic for input: %q", input)
			}
		}()
		out := ToConstant(input)
		if !utf8.ValidString(out) {
			t.Errorf("output is not valid UTF-8: %q", out)
		}
	})
}

func FuzzToHeader(f *testing.F) {
	f.Add("hello world")
	f.Add("HELLO_WORLD")
	f.Add("123-abc")
	f.Add("")
	f.Add("こんにちは 世界")
	f.Add("😀😃😄")
	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic for input: %q", input)
			}
		}()
		out := ToHeader(input)
		if !utf8.ValidString(out) {
			t.Errorf("output is not valid UTF-8: %q", out)
		}
	})
}

func FuzzToSentence(f *testing.F) {
	f.Add("hello world")
	f.Add("HELLO_WORLD")
	f.Add("123-abc")
	f.Add("")
	f.Add("こんにちは 世界")
	f.Add("😀😃😄")
	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic for input: %q", input)
			}
		}()
		out := ToSentence(input)
		if !utf8.ValidString(out) {
			t.Errorf("output is not valid UTF-8: %q", out)
		}
	})
}

func FuzzToSnake(f *testing.F) {
	f.Add("hello world")
	f.Add("HELLO_WORLD")
	f.Add("123-abc")
	f.Add("")
	f.Add("こんにちは 世界")
	f.Add("😀😃😄")
	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic for input: %q", input)
			}
		}()
		out := ToSnake(input)
		if !utf8.ValidString(out) {
			t.Errorf("output is not valid UTF-8: %q", out)
		}
	})
}

func FuzzToKebab(f *testing.F) {
	f.Add("hello world")
	f.Add("HELLO_WORLD")
	f.Add("123-abc")
	f.Add("")
	f.Add("こんにちは 世界")
	f.Add("😀😃😄")
	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("panic for input: %q", input)
			}
		}()
		out := ToKebab(input)
		if !utf8.ValidString(out) {
			t.Errorf("output is not valid UTF-8: %q", out)
		}
	})
}
