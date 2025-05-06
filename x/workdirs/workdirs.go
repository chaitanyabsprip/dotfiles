package workdirs

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"github.com/charlievieth/fastwalk"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

// Workdirs searches each path in paths for git repositories
// and returns the absolute paths or the repositories. Additionally it
// also looks for the PROJECTS environment variable and returns all
// directories in that path.
func Workdirs() []string {
	workdirs := parallelFindDirsIn(
		env.XDG_CONFIG_HOME,
		env.PROJECTS,
		env.DOTFILES,
	)
	workdirs = append(
		workdirs,
		parallelFindGitDirs(env.PROJECTS, env.PROGRAMS)...)
	workdirs = append(
		workdirs,
		env.PROJECTS,
		env.PROGRAMS,
		env.SCRIPTS,
		env.DOTFILES,
		env.NOTESPATH,
		env.DOWNLOADS,
		env.PICTURES,
	)
	return dedupe(workdirs)
}

func parallelFindGitDirs(paths ...string) []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allDirs []string

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			dirs := findGitDirs(p)
			mu.Lock()
			allDirs = append(allDirs, dirs...)
			mu.Unlock()
		}(path)
	}
	wg.Wait()
	return allDirs
}

func parallelFindDirsIn(paths ...string) []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allDirs []string

	for _, path := range paths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			dirs := findDirsIn(p)
			mu.Lock()
			allDirs = append(allDirs, dirs...)
			mu.Unlock()
		}(path)
	}
	wg.Wait()
	return allDirs
}

// findDirsIn finds all directories in path for depth 1 only.
func findDirsIn(path string) []string {
	if path == "" {
		return nil
	}
	dirs := make([]string, 0)
	fsdirs, err := os.ReadDir(path)
	if err != nil {
		return nil
	}
	for _, fsdir := range fsdirs {
		name := fsdir.Name()
		if _, skip := skipList[name]; skip || name == ".git" {
			continue
		}
		if fsdir.IsDir() {
			dirs = append(dirs, filepath.Join(path, name))
		} else if fsdir.Type()&os.ModeSymlink != 0 {
			resolvedPath, err := resolveSymlink(
				filepath.Join(path, name),
			)
			if err != nil {
				continue
			}
			dirs = append(dirs, resolvedPath)
		}
	}
	return dirs
}

var skipList = map[string]struct{}{
	"node_modules": {},
	"flutter":      {},
	".venv":        {},
	"nvm":          {},
	".terraform":   {},
}

func findGitDirs(path string) []string {
	var dirs []string
	mu := sync.Mutex{}
	baseDepth := strings.Count(path, string(os.PathSeparator))
	maxDepth := 6

	err := fastwalk.Walk(
		&fastwalk.DefaultConfig,
		path,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			currDepth := strings.Count(path, string(os.PathSeparator))
			if currDepth-baseDepth > maxDepth {
				return filepath.SkipDir
			}
			name := d.Name()
			if _, skip := skipList[name]; skip {
				return filepath.SkipDir
			}
			if name == ".git" {
				mu.Lock()
				dirs = append(dirs, filepath.Dir(path))
				mu.Unlock()
				return filepath.SkipDir
			}
			return nil
		},
	)
	if err != nil {
		return nil
	}
	return dirs
}

func dedupe(slice []string) []string {
	seen := make(map[string]struct{}, len(slice))
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if _, yes := seen[item]; !yes {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// resolveSymlink resolves any symlinks in path. It returns the absolute
// path to the resolved symlink. If the resolved symlink is not a
// directory then an error is returned.
func resolveSymlink(path string) (string, error) {
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	var absolutePath string
	if filepath.IsAbs(resolvedPath) {
		absolutePath = resolvedPath
	} else {
		absolutePath = filepath.Join(
			filepath.Dir(path), resolvedPath,
		)
	}
	if info, err := os.Stat(absolutePath); err == nil &&
		info.IsDir() {
		return absolutePath, nil
	}
	return "", fmt.Errorf("%s is not a directory", absolutePath)
}

// Worktrees searches "$PROJECTS" path for git worktrees and returns the
// absolute paths of the worktrees.
func Worktrees() []string {
	worktrees := make([]string, 0)
	mu := sync.Mutex{}
	semaphore := make(chan struct{}, 8)

	walkfn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if _, skip := skipList[name]; skip {
			return filepath.SkipDir
		}
		if name == ".git" {
			if d.IsDir() || d.Type()&os.ModeSymlink != 0 {
				return filepath.SkipDir
			}
			semaphore <- struct{}{}
			go func(gitDir string) {
				defer func() { <-semaphore }()
				if isWorktree(gitDir) && !isSubmodule(gitDir) {
					mu.Lock()
					worktrees = append(worktrees, gitDir)
					mu.Unlock()
				}
			}(filepath.Dir(path))
			return filepath.SkipDir
		}
		return nil
	}

	err := fastwalk.Walk(
		&fastwalk.DefaultConfig,
		env.PROJECTS,
		walkfn,
	)
	if err != nil {
		return nil
	}

	return worktrees
}

func isWorktree(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = path
	output, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(output)) == "true"
}

func isSubmodule(path string) bool {
	cmd := exec.Command(
		"git",
		"rev-parse",
		"--show-superproject-working-tree",
	)
	cmd.Dir = path
	output, err := cmd.Output()
	return err == nil && len(strings.TrimSpace(string(output))) > 0
}

func Shorten(paths []string) []string {
	shortPaths := make([]string, 0)
	for _, path := range paths {
		sPath := strings.TrimPrefix(path, env.HOME)
		sPath = strings.TrimLeftFunc(sPath, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		shortPaths = append(shortPaths, sPath)
	}
	return shortPaths
}
