// Package embed provides utilities for working with embedded files and directories.
// It facilitates copying embedded filesystem contents to the host system,
// handling operations such as setup, configuration deployment, and file management
// with appropriate permissions and path handling.
package embed

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aymanbagabas/go-udiff"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/manifest"
	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
	"github.com/Chaitanyabsprip/dotfiles/x/have"
)

const Skip = ""

// CheckMode is how copy behaves when Check is set to something other than
// CheckOff — used by `dot status`/`dot diff` (internal/dot) to preview
// drift across tools without writing anything.
type CheckMode int

const (
	CheckOff  CheckMode = iota // normal deploy (the default)
	CheckList                  // record drifted paths in Drifted, write nothing
	CheckDiff                  // CheckList, and also print each drifted file's diff
)

// Check selects copy's mode for the whole process — see CheckMode. Callers
// must reset it to CheckOff when done; it's a package var (not threaded
// through every call) because SetupCmds are already composed and run via
// `cmd.Run()` with no channel to pass a mode through.
var Check CheckMode

// Drifted collects dest paths found to have local changes during a
// CheckList/CheckDiff run. Callers reset it (Drifted = nil) before each
// run they want isolated to one tool.
var Drifted []string

// SetupAll deploys every file embedded under name into configDir,
// per-file drift-aware (see copy). It replaces an older, cruder
// implementation that deleted the tool's whole config directory before
// every deploy — which would have destroyed any local edit before drift
// detection ever got a chance to see it.
func SetupAll(
	embedFs embed.FS,
	name, configDir string,
	overrides map[string]string,
) error {
	return CopyAllFiles(embedFs, name, configDir, overrides)
}

func CopyFilesRegx(
	embedFs embed.FS,
	name, configDir, pattern string,
	overrides map[string]string,
) error {
	if len(configDir) == 0 {
		configDir = filepath.Join(env.Home, ".config")
	}
	regx, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	m, err := manifest.Load()
	if err != nil {
		return err
	}
	walkErr := fs.WalkDir(
		embedFs,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if regx.Match([]byte(path)) {
				targetPath := filepath.Join(
					configDir,
					path,
				)
				if altPath, ok := overrides[path]; ok {
					if altPath == Skip {
						return nil
					}
					targetPath = altPath
				}
				return copy(embedFs, m, d, path, targetPath)
			}
			return nil
		},
	)
	if walkErr != nil {
		return walkErr
	}
	if Check != CheckOff {
		return nil
	}
	return m.Save()
}

func CopyAllFiles(
	embedFs embed.FS,
	name, configDir string,
	overrides map[string]string,
) error {
	if len(configDir) == 0 {
		configDir = filepath.Join(env.Home, ".config")
	}
	m, err := manifest.Load()
	if err != nil {
		return err
	}
	walkErr := fs.WalkDir(
		embedFs,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			targetPath := filepath.Join(configDir, path)
			if altPath, ok := overrides[path]; ok {
				if altPath == Skip {
					return nil
				}
				targetPath = altPath
			}
			return copy(embedFs, m, d, path, targetPath)
		},
	)
	if walkErr != nil {
		return walkErr
	}
	if Check != CheckOff {
		return nil
	}
	return m.Save()
}

// copy deploys a single embedded path to dest. For a regular file it is
// drift-aware: a dest that has been hand-edited since dot last deployed
// it is left alone (diff printed, manifest untouched) unless DOT_FORCE is
// set, matching VISION.md's three-way comparison of live vs. manifest vs.
// embedded content. A dest that matches what the manifest last recorded
// is always safe to update — that's just picking up an upstream config
// change, not overwriting a user edit.
func copy(embedFs embed.FS, m manifest.Manifest, d fs.DirEntry, path, dest string) error {
	if d.IsDir() {
		if Check != CheckOff {
			return nil
		}
		return os.MkdirAll(dest, 0o755)
	}
	content, err := fs.ReadFile(embedFs, path)
	if err != nil {
		return err
	}
	live, err := os.ReadFile(dest)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	drifted := false
	if err == nil { // dest already exists — check for drift
		liveHash := manifest.Hash(live)
		drifted = liveHash != manifest.Hash(content) && liveHash != m[dest]
	}

	if Check != CheckOff {
		if drifted {
			Drifted = append(Drifted, dest)
			if Check == CheckDiff {
				printDiff(dest, live, content)
			}
		}
		return nil
	}

	if drifted && os.Getenv(`DOT_FORCE`) == `` {
		fmt.Printf("drift: %s has local changes, skipping (set DOT_FORCE=1 to overwrite)\n", dest)
		printDiff(dest, live, content)
		return nil
	}
	if err := os.WriteFile(dest, content, getFileMode(path)); err != nil {
		return err
	}
	m[dest] = manifest.Hash(content)
	return nil
}

// prettyDiffTools are external diff pagers tried in order, each fed the
// plain unified diff on stdin. First one found on PATH wins; if none are
// installed, or the one found fails, printDiff falls back to the plain
// unified diff.
var prettyDiffTools = []struct {
	name string
	args []string
}{
	{`delta`, []string{`--paging=never`}},
	{`diff-so-fancy`, nil},
}

func printDiff(dest string, live, embedded []byte) {
	unified := udiff.Unified(dest+` (local)`, dest+` (dot)`, string(live), string(embedded))
	for _, tool := range prettyDiffTools {
		if ok, _ := have.Executable(tool.name); !ok {
			continue
		}
		cmd := exec.Command(tool.name, tool.args...)
		cmd.Stdin = strings.NewReader(unified)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return
		}
		break
	}
	fmt.Println(unified)
}

func getFileMode(path string) fs.FileMode {
	var mode fs.FileMode
	if strings.Contains(path, "/bin/") ||
		strings.HasPrefix(path, "bin/") {
		mode = 0o755
	} else {
		mode = 0o644
	}
	return mode
}
