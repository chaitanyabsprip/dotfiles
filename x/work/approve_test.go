package work

import (
	"errors"
	"testing"
)

func TestApproveRejectsNonPRURL(t *testing.T) {
	tests := []string{
		``,
		`123`,
		`https://github.com/owner/repo`,
		`https://github.com/owner/repo/issues/123`,
		`https://gitlab.com/owner/repo/pull/123`,
		`not a url`,
	}
	for _, in := range tests {
		if err := approve(in); err == nil {
			t.Errorf("approve(%q) should reject a non-PR-URL", in)
		}
	}
}

func TestApproveCallsGhForValidURL(t *testing.T) {
	orig := ghApprove
	defer func() { ghApprove = orig }()

	var got string
	ghApprove = func(url string) error { got = url; return nil }

	url := `https://github.com/owner/repo/pull/123`
	if err := approve(url); err != nil {
		t.Fatal(err)
	}
	if got != url {
		t.Errorf("ghApprove called with %q, want %q", got, url)
	}
}

func TestApproveAcceptsTrailingSlash(t *testing.T) {
	orig := ghApprove
	defer func() { ghApprove = orig }()
	ghApprove = func(string) error { return nil }

	if err := approve(`https://github.com/owner/repo/pull/123/`); err != nil {
		t.Errorf("trailing slash should be accepted: %v", err)
	}
}

func TestApprovePropagatesGhError(t *testing.T) {
	orig := ghApprove
	defer func() { ghApprove = orig }()
	wantErr := errors.New(`gh failed`)
	ghApprove = func(string) error { return wantErr }

	if err := approve(`https://github.com/owner/repo/pull/1`); !errors.Is(err, wantErr) {
		t.Errorf("approve error = %v, want %v", err, wantErr)
	}
}
