package workdirs

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charlievieth/fastwalk"
)

// Workdirs searches each path in paths for git repositories
// and returns the absolute paths or the repositories. Additionally it
// also looks for the PROJECTS environment variable and returns all
// directories in that path.
func Workdirs() []string {
	workdirs := make([]string, 0)
	workdirs = append(
		workdirs,
		findDirsIn(os.Getenv("XDG_CONFIG_HOME"))...)
	workdirs = append(
		workdirs,
		findDirsIn(os.Getenv("PROJECTS"))...,
	)
	workdirs = append(
		workdirs,
		findDirsIn(os.Getenv("DOTFILES"))...,
	)
	workdirs = append(
		workdirs,
		findGitDirs(os.Getenv("PROJECTS"))...,
	)
	workdirs = append(
		workdirs,
		os.Getenv("PROJECTS"),
		os.Getenv("SCRIPTS"),
		os.Getenv("DOTFILES"),
		os.Getenv("NOTESPATH"),
		os.Getenv("DOWNLOADS"),
		filepath.Join(os.Getenv("HOME"), "projects"),
		filepath.Join(os.Getenv("HOME"), "programs"),
		filepath.Join(os.Getenv("HOME"), "pictures"),
	)
	return dedupe(workdirs)
}

// findDirsIn finds all directories in path for depth 1 only.
func findDirsIn(path string) []string {
	dirs := make([]string, 0)
	fsdirs, err := os.ReadDir(path)
	if err != nil {
		return []string{}
	}
	for _, fsdir := range fsdirs {
		if fsdir.IsDir() {
			dirs = append(
				dirs,
				filepath.Join(path, fsdir.Name()),
			)
		} else if fsdir.Type()&os.ModeSymlink != 0 {
			resolvedPath, err := resolveSymlink(
				filepath.Join(path, fsdir.Name()),
			)
			if err != nil {
				continue
			}
			dirs = append(dirs, resolvedPath)
		} else {
		}
	}
	return dirs
}

func findGitDirs(path string) []string {
	var dirs []string
	dirCh := make(chan string)

	go func() {
		for dir := range dirCh {
			dirs = append(dirs, dir)
		}
	}()

	err := fastwalk.Walk(
		&fastwalk.DefaultConfig,
		path,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil || !d.IsDir() {
				return nil
			}
			if d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			if d.Name() == ".git" {
				dirCh <- filepath.Dir(path)
				return filepath.SkipDir
			}
			return nil
		},
	)
	if err != nil {
		return []string{}
	}
	close(dirCh)
	return dirs
}

func dedupe(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
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
	dirCh := make(chan string)

	go func() {
		for dir := range dirCh {
			worktrees = append(worktrees, dir)
		}
	}()

	walkfn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		switch d.Name() {
		case `.venv`, `node_modules`:
			return filepath.SkipDir
		case `.git`:
			if d.IsDir() || d.Type()&os.ModeSymlink != 0 {
				return filepath.SkipDir
			}
			if isWorktree(filepath.Dir(path)) &&
				!isSubmodule(filepath.Dir(path)) {
				dirCh <- filepath.Dir(path)
				return filepath.SkipDir
			}
		}
		return nil
	}

	err := fastwalk.Walk(
		&fastwalk.DefaultConfig,
		os.Getenv("PROJECTS"),
		walkfn,
	)
	if err != nil {
		return []string{}
	}

	close(dirCh)
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
		sPath := strings.ReplaceAll(path, os.Getenv("HOME"), "")
		sPath = strings.TrimLeftFunc(sPath, func(r rune) bool {
			return !unicode.IsLetter(r)
		})
		shortPaths = append(shortPaths, sPath)
	}
	return shortPaths
}
