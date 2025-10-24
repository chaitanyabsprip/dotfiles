package caseconv

import (
	"testing"
)

func TestToLower(t *testing.T) {
	if got := ToLower("Hello World"); got != "hello world" {
		t.Errorf("ToLower failed: got %q", got)
	}
}

func TestToUpper(t *testing.T) {
	if got := ToUpper("Hello World"); got != "HELLO WORLD" {
		t.Errorf("ToUpper failed: got %q", got)
	}
}

func TestToCamel(t *testing.T) {
	if got := ToCamel("hello world"); got != "helloWorld" {
		t.Errorf("ToCamel failed: got %q", got)
	}
	if got := ToCamel("hello_world"); got != "helloWorld" {
		t.Errorf("ToCamel failed: got %q", got)
	}
	if got := ToCamel("hello-world"); got != "helloWorld" {
		t.Errorf("ToCamel failed: got %q", got)
	}
}

func TestToTitle(t *testing.T) {
	if got := ToTitle("hello world"); got != "Hello World" {
		t.Errorf("ToTitle failed: got %q", got)
	}
}

func TestToConstant(t *testing.T) {
	if got := ToConstant("hello world"); got != "HELLO_WORLD" {
		t.Errorf("ToConstant failed: got %q", got)
	}
}

func TestToHeader(t *testing.T) {
	if got := ToHeader("hello world"); got != "Hello-World" {
		t.Errorf("ToHeader failed: got %q", got)
	}
}

func TestToSentence(t *testing.T) {
	if got := ToSentence("hello world"); got != "Hello world" {
		t.Errorf("ToSentence failed: got %q", got)
	}
}

func TestToSnake(t *testing.T) {
	if got := ToSnake("hello world"); got != "hello_world" {
		t.Errorf("ToSnake failed: got %q", got)
	}
}

func TestToKebab(t *testing.T) {
	if got := ToKebab("hello world"); got != "hello-world" {
		t.Errorf("ToKebab failed: got %q", got)
	}
}
