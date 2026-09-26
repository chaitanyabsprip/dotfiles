package tmux

import (
	"path/filepath"
	"reflect"
	"testing"
)

var (
	api = bookmark{Session: `api`, Dir: `/w/api`}
	cli = bookmark{Session: `cli`, Dir: `/w/cli`, Target: `2.1`}
	web = bookmark{Session: `web`, Dir: `/w/web`}
)

func TestParseFormatBookmarks(t *testing.T) {
	got := parseBookmarks("api=/w/api\ncli=/w/cli:2.1\n\nnot a bookmark\n")
	if want := []bookmark{api, cli}; !reflect.DeepEqual(got, want) {
		t.Fatalf("parse = %+v, want %+v", got, want)
	}
	if out := formatBookmarks(got); out != "api=/w/api\ncli=/w/cli:2.1\n" {
		t.Errorf("format = %q", out)
	}
}

func TestAddBookmark(t *testing.T) {
	bs := addBookmark([]bookmark{api}, cli)
	if want := []bookmark{api, cli}; !reflect.DeepEqual(bs, want) {
		t.Errorf("add new = %+v", bs)
	}
	if again := addBookmark(bs, api); !reflect.DeepEqual(again, bs) {
		t.Errorf("add duplicate changed the list: %+v", again)
	}
}

func TestReplaceBookmark(t *testing.T) {
	tests := []struct {
		slot int
		want []bookmark
	}{
		{1, []bookmark{web, cli}},      // replace first
		{2, []bookmark{api, web}},      // replace last
		{3, []bookmark{api, cli, web}}, // past the end appends
	}
	for _, tt := range tests {
		got, err := replaceBookmark([]bookmark{api, cli}, tt.slot, web)
		if err != nil || !reflect.DeepEqual(got, tt.want) {
			t.Errorf("slot %d: got %+v, %v; want %+v", tt.slot, got, err, tt.want)
		}
	}
	got, _ := replaceBookmark([]bookmark{api, cli}, 2, api) // would duplicate api
	if want := []bookmark{api}; !reflect.DeepEqual(got, want) {
		t.Errorf("dedupe: got %+v", got)
	}
	if _, err := replaceBookmark([]bookmark{api}, 0, web); err == nil {
		t.Error(`slot 0 should fail`)
	}
}

func TestRemoveBookmarks(t *testing.T) {
	api2 := bookmark{Session: `api2`, Dir: `/w/api2`}
	apiPane := bookmark{Session: `api`, Dir: `/w/api`, Target: `1.2`}
	got := removeBookmarks([]bookmark{api, api2, apiPane}, `api`)
	if want := []bookmark{api2}; !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// The standalone harpoon reads the same file; the locations must agree.
func TestHarpoonFileMatchesStandalone(t *testing.T) {
	tests := []struct{ xdg, goos, want string }{
		{`/x/cache`, `darwin`, `/x/cache/.tmux-harpoon-sessions`},
		{``, `darwin`, `/h/Library/Caches/.tmux-harpoon-sessions`},
		{``, `linux`, `/h/.cache/.tmux-harpoon-sessions`},
	}
	for _, tt := range tests {
		if got := harpoonFileFor(tt.xdg, tt.goos, `/h`); got != filepath.FromSlash(tt.want) {
			t.Errorf("harpoonFileFor(%q, %q) = %q, want %q", tt.xdg, tt.goos, got, tt.want)
		}
	}
}

func TestLoadSaveBookmarks(t *testing.T) {
	path := filepath.Join(t.TempDir(), `sub`, `harpoon`)
	bs, err := loadBookmarks(path)
	if err != nil || bs != nil {
		t.Fatalf("missing file: %+v, %v", bs, err)
	}
	if err := saveBookmarks(path, []bookmark{api, cli}); err != nil {
		t.Fatal(err)
	}
	bs, err = loadBookmarks(path)
	if err != nil || !reflect.DeepEqual(bs, []bookmark{api, cli}) {
		t.Errorf("round trip: %+v, %v", bs, err)
	}
}
