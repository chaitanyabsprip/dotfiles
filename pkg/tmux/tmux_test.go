package tmux

import "testing"

// TestFindSessionExactNameMatch guards against matching session names
// by raw substring: querying for "ss" must not match an existing
// "ss_ui" session just because "ss_ui" starts with "ss".
func TestFindSessionExactNameMatch(t *testing.T) {
	sessions := []Session{
		{Name: `dotfiles`, Path: `/Users/x/projects/dotfiles`},
		{Name: `ss_ui`, Path: `/Users/x/projects/solar-square/ss_ui`},
	}

	if name, _ := findSession(sessions, Session{Name: `ss`}); name != `` {
		t.Errorf("findSession(Name: ss) = %q, want no match", name)
	}

	if name, path := findSession(sessions, Session{Name: `ss_ui`}); name != `ss_ui` ||
		path != `/Users/x/projects/solar-square/ss_ui` {
		t.Errorf("findSession(Name: ss_ui) = %q, %q, want exact match", name, path)
	}
}

func TestFindSessionExactPathMatch(t *testing.T) {
	sessions := []Session{
		{Name: `ss_ui`, Path: `/Users/x/projects/solar-square/ss_ui`},
	}

	if name, _ := findSession(
		sessions, Session{Path: `/Users/x/projects/chaitanyasharma-sse/ss`},
	); name != `` {
		t.Errorf("findSession(Path: .../ss) = %q, want no match", name)
	}
}
