// Package with provides utilities for temporarily modifying process state such as
// working directory, environment variables, and PATH, with automatic restoration.
package with

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

// PopFunc is a function that restores a previous process state, such as
// working directory or environment variable, and returns an error if restoration fails.
type PopFunc func() error

// Pwd changes the current working directory to dir and returns a PopFunc that
// restores the previous working directory. If changing the directory fails,
// it returns an error.
func Pwd(dir string) (pop PopFunc, err error) {
	pwd := ``
	if pwd, err = os.Getwd(); err != nil {
		return nil, err
	}
	if err = os.Chdir(dir); err != nil {
		return nil, err
	}
	return func() error { return os.Chdir(pwd) }, nil
}

// Env sets the environment variable name to value and returns a PopFunc that
// restores the previous value (or unsets it if it did not exist). If setting
// the variable fails, it returns an error.
func Env(name, value string) (PopFunc, error) {
	old, exists := os.LookupEnv(name)
	if err := os.Setenv(name, value); err != nil {
		return nil, err
	}
	if exists {
		return func() error { return os.Setenv(name, old) }, nil
	} else {
		return func() error { return os.Unsetenv(name) }, nil
	}
}

// Path prepends dir to the PATH environment variable and returns a PopFunc that
// restores the previous PATH value. If setting the variable fails, it returns an error.
func Path(dir string) (PopFunc, error) {
	newPath := fmt.Sprintf(
		`%s%c%s`,
		dir, filepath.ListSeparator, env.Path,
	)
	return Env(`PATH`, newPath)
}
