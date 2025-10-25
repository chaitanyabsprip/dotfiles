package base64

import (
	"testing"
)

func TestEncodeDecodeStandard(t *testing.T) {
	input := "hello world"
	encoded, err := Encode(input, false)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	if encoded != "aGVsbG8gd29ybGQ=" {
		t.Errorf("Expected aGVsbG8gd29ybGQ=, got %s", encoded)
	}
	decoded, err := Decode(encoded, false)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if decoded != input {
		t.Errorf("Expected %s, got %s", input, decoded)
	}
}

func TestEncodeDecodeURLSafe(t *testing.T) {
	input := "hello world?"
	encoded, err := Encode(input, true)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	if encoded != "aGVsbG8gd29ybGQ_" {
		t.Errorf("Expected aGVsbG8gd29ybGQ_, got %s", encoded)
	}
	decoded, err := Decode(encoded, true)
	if err != nil {
		t.Fatalf("Decode error: %v", err)
	}
	if decoded != input {
		t.Errorf("Expected %s, got %s", input, decoded)
	}
}

func TestDecodeInvalid(t *testing.T) {
	_, err := Decode("not_base64!!", false)
	if err == nil {
		t.Error("Expected error for invalid base64 input")
	}
}
