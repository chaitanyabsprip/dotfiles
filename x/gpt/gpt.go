package gpt

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

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
			or(opts.Model, ``),
		)
	}
	argList := []string{`mods`}
	if len(opts.Format) > 0 {
		argList = append(argList, fmt.Sprintf(`-f %s`, opts.Format))
	}
	if len(opts.Model) > 0 {
		argList = append(argList, fmt.Sprintf(`-m=%s`, opts.Model))
	}
	if opts.Quiet {
		argList = append(argList, `-q`)
	}
	if len(opts.Role) > 0 {
		argList = append(argList, fmt.Sprintf(`--role=%s`, opts.Role))
	}
	if len(opts.StatusText) > 0 {
		argList = append(
			argList,
			fmt.Sprintf(`--status-text=%s`, opts.StatusText),
		)
	}
	if opts.NoCache {
		argList = append(argList, `--no-cache`)
	} else {
		if _, err := FindConversation(opts.Title); err != nil {
			argList = append(argList, fmt.Sprintf(`-t=%s`, opts.Title))
		} else {
			argList = append(argList, fmt.Sprintf(`-c=%s`, opts.Title))
		}
	}
	return append(argList, fmt.Sprintf(`"%s"`, opts.Query))
}

func Exec(opts GptOpts) error {
	argList := BuildCmdStr(opts)
	os.Setenv(`CLICOLOR_FORCE`, `1`)
	// this defer won't run, but is here for completeness
	defer os.Unsetenv(`CLICOLOR_FORCE`)
	return run.SysExec(argList...)
}

// Run executes the gpt command and returns the output as a string.
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

func FindConversation(query string) (string, error) {
	conversations, err := ListConversations()
	if err != nil {
		return "", err
	}
	for _, conversation := range conversations {
		if strings.Contains(conversation, query) {
			return conversation, nil
		}
	}
	return "", fmt.Errorf("conversation not found: %s", query)
}

func ListConversations() ([]string, error) {
	out := run.Out(`mods`, `-l`)
	if len(out) == 0 {
		return nil, nil
	}

	lines := strings.Split(string(out), "\n")
	conversations := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		conversations = append(conversations, fields[1])
	}
	return conversations, nil
}

func or(a, b string) string {
	if len(a) > 0 {
		return a
	}
	return b
}
