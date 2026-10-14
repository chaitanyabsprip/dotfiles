package json

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUnescape(t *testing.T) {
	in := `[{"key":"k","value":{"op":"c","payload":"{\"id\":12345678901234567890,\"tags\":\"[1,2]\"}"},` +
		`"headers":{"n":"42","t":"true","p":"plain","bad":"{nope"}}]`
	out, err := Unescape([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	var got, want any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	json.Unmarshal([]byte(`[{"key":"k","value":{"op":"c","payload":{"id":12345678901234567890,"tags":[1,2]}},`+
		`"headers":{"n":"42","t":"true","p":"plain","bad":"{nope"}}]`), &want)
	if g, w := mustJSON(got), mustJSON(want); g != w {
		t.Errorf("got  %s\nwant %s", g, w)
	}
	if !strings.Contains(string(out), "12345678901234567890") {
		t.Errorf("big id lost precision: %s", out)
	}
}

func TestUnescapeTopLevelString(t *testing.T) {
	out, err := Unescape([]byte(`"{\"a\":1}"`))
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "{\n  \"a\": 1\n}" {
		t.Errorf("got %s", out)
	}
}

func TestUnescapeInvalid(t *testing.T) {
	if _, err := Unescape([]byte(`{nope`)); err == nil {
		t.Error("expected error for invalid input")
	}
}

func mustJSON(v any) string { b, _ := json.Marshal(v); return string(b) }
