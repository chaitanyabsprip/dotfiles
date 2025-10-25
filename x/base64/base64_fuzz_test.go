package base64

import (
	"testing"
)

func FuzzDecode(f *testing.F) {
	f.Add("aGVsbG8gd29ybGQ=") // valid base64
	f.Add("aGVsbG8gd29ybGQ_") // valid url-safe base64
	f.Add("not_base64!!")     // invalid

	f.Fuzz(func(t *testing.T, input string) {
		for _, urlSafe := range []bool{false, true} {
			_, _ = Decode(input, urlSafe)
		}
	})
}
