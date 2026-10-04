package work

import (
	"errors"
	"testing"
)

func TestApproveRejectsNonPRURL(t *testing.T) {
	orig := ghRepoURL
	defer func() { ghRepoURL = orig }()
	ghRepoURL = func() (string, error) { return ``, errors.New(`not a git repo`) }

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

func TestApproveExpandsPRNumber(t *testing.T) {
	origApprove, origRepo := ghApprove, ghRepoURL
	defer func() { ghApprove, ghRepoURL = origApprove, origRepo }()

	var got string
	ghApprove = func(url string) error { got = url; return nil }
	ghRepoURL = func() (string, error) { return `https://github.com/owner/repo`, nil }

	if err := approve(`42`); err != nil {
		t.Fatal(err)
	}
	if want := `https://github.com/owner/repo/pull/42`; got != want {
		t.Errorf("ghApprove called with %q, want %q", got, want)
	}
}

func TestApproveRejectsPRNumberOutsideGitHub(t *testing.T) {
	origApprove, origRepo := ghApprove, ghRepoURL
	defer func() { ghApprove, ghRepoURL = origApprove, origRepo }()

	ghApprove = func(string) error { t.Error(`ghApprove should not be called`); return nil }
	ghRepoURL = func() (string, error) { return `https://gitlab.com/owner/repo`, nil }

	if err := approve(`42`); err == nil {
		t.Error(`a non-GitHub repo should be rejected`)
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
