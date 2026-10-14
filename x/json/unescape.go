package json

import (
	"bytes"
	"encoding/json"
	"strings"
)

// Unescape parses data and replaces every string that holds a JSON object
// or array with the parsed value, recursively, so JSON stringified inside
// JSON (a payload inside an envelope) comes out as one nested document.
// Other strings, including "42" or "true", are left alone.
func Unescape(data []byte) ([]byte, error) {
	v, err := decode(data)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(expand(v), "", "  ")
}

func expand(v any) any {
	switch t := v.(type) {
	case string:
		s := strings.TrimSpace(t)
		if !strings.HasPrefix(s, "{") && !strings.HasPrefix(s, "[") {
			return t
		}
		inner, err := decode([]byte(s))
		if err != nil {
			return t
		}
		return expand(inner)
	case map[string]any:
		for k, e := range t {
			t[k] = expand(e)
		}
	case []any:
		for i, e := range t {
			t[i] = expand(e)
		}
	}
	return v
}

// decode keeps numbers as json.Number so large IDs don't lose precision
// through float64.
func decode(data []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
