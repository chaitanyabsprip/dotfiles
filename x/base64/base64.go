// Package base64 provides functions to encode and decode strings using Base64 encoding.
package base64

import (
	"encoding/base64"
	"errors"
)

func Encode(input string, urlSafe bool) (string, error) {
	var enc *base64.Encoding
	if urlSafe {
		enc = base64.URLEncoding
	} else {
		enc = base64.StdEncoding
	}
	return enc.EncodeToString([]byte(input)), nil
}

func Decode(input string, urlSafe bool) (string, error) {
	var enc *base64.Encoding
	if urlSafe {
		enc = base64.URLEncoding
	} else {
		enc = base64.StdEncoding
	}
	decoded, err := enc.DecodeString(input)
	if err != nil {
		return "", errors.New("invalid base64 input")
	}
	return string(decoded), nil
}
