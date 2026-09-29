package work

import (
	"strings"
	"testing"

	"github.com/arl/gitstatus"
)

func TestSummarizeClean(t *testing.T) {
	if got := summarize(&gitstatus.Status{IsClean: true}); got != "✔" {
		t.Errorf("summarize(clean) = %q, want ✔", got)
	}
}

func TestSummarizeCleanWithStash(t *testing.T) {
	s := &gitstatus.Status{IsClean: true, NumStashed: 2}
	got := summarize(s)
	if want := " 2 ✔"; got != want {
		t.Errorf("summarize(clean, stashed) = %q, want %q", got, want)
	}
}

func TestSummarizeDirty(t *testing.T) {
	s := &gitstatus.Status{
		Insertions: 6, Deletions: 7,
		Porcelain: gitstatus.Porcelain{
			NumStaged: 2, NumModified: 1, NumConflicts: 3, NumUntracked: 4,
			AheadCount: 1, BehindCount: 5,
		},
	}
	got := summarize(s)
	for _, want := range []string{
		" 6", " 7", " 2", `! 3`, " 1", `?? 4`,
		"\U000f0dbc 1", "\U000f0db9 5",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("summarize(dirty) = %q, missing %q", got, want)
		}
	}
	if strings.Contains(got, `clean`) || strings.Contains(got, `staged`) {
		t.Errorf("summarize(dirty) = %q, should use symbols not words", got)
	}
}

func TestSummarizeOmitsZeroCounts(t *testing.T) {
	s := &gitstatus.Status{Porcelain: gitstatus.Porcelain{NumModified: 1}}
	got := summarize(s)
	if want := " 1"; got != want {
		t.Errorf("summarize(only modified) = %q, want %q", got, want)
	}
}

func TestPRSymbol(t *testing.T) {
	tests := map[string]string{
		`MERGED`:  `✓`,
		`CLOSED`:  `✗`,
		`OPEN`:    `●`,
		`UNKNOWN`: `?`,
	}
	for state, want := range tests {
		if got := prSymbol(state); got != want {
			t.Errorf("prSymbol(%q) = %q, want %q", state, got, want)
		}
	}
}
