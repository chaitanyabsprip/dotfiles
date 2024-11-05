package gpt

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/rwxrob/bonzai/run"
)

type GptOpts struct {
	Model      string
	Format     string
	Query      string
	Role       string
	StatusText string
	Title      string
	NoCache    bool
	Quiet      bool
	Stdin      io.Reader
}

func BuildCmdStr(opts GptOpts) []string {
	if len(opts.Title) == 0 {
		opts.Title = fmt.Sprintf(
			`%s:%s`,
			or(opts.Role, "default"),
			or(opts.Model, "liquid/lfm-40b"),
		)
	}
	argList := []string{`mods`}
	if len(opts.Model) > 0 {
		argList = append(argList, `-m`, opts.Model)
	}
	if len(opts.Role) > 0 {
		argList = append(argList, `--role`, opts.Role)
	}
	if len(opts.Format) > 0 {
		argList = append(argList, `-f`, opts.Format)
	}
	if opts.Quiet {
		argList = append(argList, `-q`)
	}
	if opts.NoCache {
		argList = append(argList, `--no-cache`)
	} else {
		argList = append(argList, `-t`, opts.Title, `-c`, opts.Title)
	}
	if len(opts.StatusText) > 0 {
		argList = append(
			argList,
			`--status-text`,
			opts.StatusText,
		)
	}
	return append(argList, opts.Query)
}

func Exec(opts GptOpts) error {
	argList := BuildCmdStr(opts)
	os.Setenv(`CLICOLOR_FORCE`, `1`)
	// this defer won't run, but is here for completeness
	defer os.Unsetenv(`CLICOLOR_FORCE`)
	return run.SysExec(argList...)
}

func Run(opts GptOpts) (string, error) {
	argList := BuildCmdStr(opts)
	os.Setenv(`CLICOLOR_FORCE`, `1`)
	defer os.Unsetenv(`CLICOLOR_FORCE`)
	if err := run.SysExec(argList...); err != nil {
		return "", err
	}
	cmd := exec.Command(argList[0], argList[1:]...)
	cmd.Stdin = opts.Stdin
	cmd.Env = append(os.Environ(), `CLICOLOR_FORCE=1`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	if len(out) > 0 {
		return string(out), nil
	}
	return "", nil
}

func or(a, b string) string {
	if len(a) > 0 {
		return a
	}
	return b
}
