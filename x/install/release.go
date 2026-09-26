package install

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rwxrob/bonzai/github"
)

// GhRelease installs bin into BinDir from the latest GitHub release of
// repo. asset names the release file for a tag and the running GOOS and
// GOARCH, or returns "" when the project ships no build for it. A
// .tar.gz asset must hold bin at its top level (a single leading
// directory is stripped); any other asset is the binary itself.
func GhRelease(repo, bin string, asset func(tag, goos, goarch string) string) error {
	tag, err := latestTag(repo)
	if err != nil {
		return err
	}
	name := asset(tag, runtime.GOOS, runtime.GOARCH)
	if name == `` {
		return fmt.Errorf(`%s: no release build for %s/%s`, repo, runtime.GOOS, runtime.GOARCH)
	}
	if err := os.MkdirAll(BinDir, 0o755); err != nil {
		return err
	}
	dl, err := GhDownload(repo, tag, name)
	if err != nil {
		return err
	}
	defer os.Remove(dl)
	src := dl
	if strings.HasSuffix(name, `.tar.gz`) {
		// Extract inside BinDir so the final rename stays on one filesystem.
		tmp, err := os.MkdirTemp(BinDir, `.`+bin+`-`)
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmp)
		if err := ExtractTarGz(dl, tmp); err != nil {
			return err
		}
		src = filepath.Join(tmp, bin)
	}
	dest := filepath.Join(BinDir, bin)
	if err := os.Rename(src, dest); err != nil {
		return err
	}
	return os.Chmod(dest, 0o755)
}

// latestTag returns the git tag of repo's latest release. bonzai's
// github.Latest returns the release's display name, which is not always
// the tag (fd names v10.5.0 "10.5.0").
func latestTag(repo string) (string, error) {
	url := fmt.Sprintf(`https://api.%s/repos/%s/releases/latest`, github.Host, repo)
	resp, err := http.Get(url)
	if err != nil {
		return ``, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ``, fmt.Errorf(`%s: latest release: %s`, repo, resp.Status)
	}
	var rel struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return ``, err
	}
	if rel.TagName == `` {
		return ``, fmt.Errorf(`%s: latest release has no tag`, repo)
	}
	return rel.TagName, nil
}
