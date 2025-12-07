package tmux

import "testing"

func TestReconcileSessionName(t *testing.T) {
	tests := []struct {
		oldPath, newPath string
		want1, want2     string
	}{
		{"b", "b", "old-b", "b"},
		{"/a/b", "/a/b", "old-b", "b"},
		{"/a/b/c", "/a/b/d", "c", "d"},
		{"/a/b", "/a/b/c", "b", "c"},
		{"/a/b", "/a/b/c/d", "b", "d"},
		{"/x/y", "/z/w", "y", "w"},
		{"/a/b/c", "/a/d/c", "b/c", "c"},
		{"/a/b/c/d", "/a/z/c/d", "b/c/d", "d"},
		{"/a/b/c/d/e/f", "/a/b/g/d/e/f", "c/d/e/f", "f"},
		{"/a/b/c", "/a/b/d/e", "c", "e"},
		{"/a/b/root", "/a/c/root", "b/root", "c/root"},
		{"/a/b/root", "/a/b/root", "old-b/root", "b/root"},
		{"/a/b", "/a/b/c/b", "a/b", "b"},
		{"/z/a/b", "/z/a/b/c/b", "a/b", "b"},
	}
	for _, tt := range tests {
		got1, got2 := reconcileSessionName(tt.oldPath, tt.newPath)
		if got1 != tt.want1 || got2 != tt.want2 {
			t.Errorf(
				"reconcileSessionName(%q, %q) => %q, %q; want %q, %q",
				tt.oldPath,
				tt.newPath,
				got1,
				got2,
				tt.want1,
				tt.want2,
			)
		}
	}
}

// func TestCalculateCommonPath(t *testing.T) {
// 	tests := []struct {
// 		p1, p2 string
// 		want1  string
// 	}{
// 		{"/a/b", "/a/b", "a/b"},
// 		{"/a/b/c", "/a/b/d", "a/b"},
// 		{"/a/b", "/a/b/c", "a/b"},
// 		{"/x/y", "/z/w", ""},
// 	}
// 	for _, tt := range tests {
// 		got1 := calculateCommonPath(tt.p1, tt.p2)
// 		if got1 != tt.want1 {
// 			t.Errorf(
// 				"calculateCommonPath(%q, %q) => %q; want %q",
// 				tt.p1,
// 				tt.p2,
// 				got1,
// 				tt.want1,
// 			)
// 		}
// 	}
// }
