// gpt provides a user-friendly CLI for interacting with
// charmbracelet/mods CLI. It is a stateful program such that the users
// will be talking in the same conversation with consecutive calls.
package gpt

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/fn/each"
	"github.com/rwxrob/bonzai/vars"
	"github.com/rwxrob/bonzai/yq"

	"github.com/Chaitanyabsprip/dotfiles/internal/core/oscfg"
	"github.com/Chaitanyabsprip/dotfiles/x/depends"
)

const (
	ModelEnv       = `GPT_MODEL`
	RoleEnv        = `GPT_ROLE`
	TitleEnv       = `GPT_CHAT_TITLE`
	GlobalQuietEnv = `QUIET`
	QuietEnv       = `GPT_QUIET`
	StatusTextEnv  = `GPT_STATUS_TEXT`
	NoCacheEnv     = `GPT_NOCACHE`
)

var (
	modsConfPath = filepath.Join(oscfg.ConfigDir(), "mods")
	defaultModel string
)

func init() {
	var err error
	defaultModel, err = yq.EvaluateToString(
		`.default-model`,
		modsConfPath,
	)
	if err != nil {
		defaultModel = `gemini-free`
	}
}

var Cmd = &bonzai.Cmd{
	Name:  `gpt`,
	Vers:  `v0.1.0`,
	Short: `Persistent conversations with LLM models using mods`,
	Comp:  comp.Combine{comp.Cmds},
	Cmds: []*bonzai.Cmd{
		vars.Cmd,
		commitCmd,
		devCmd,
		shellCmd,
		commentCmd,
		listCmd,
	},
	Do: func(x *bonzai.Cmd, args ...string) error {
		depends.On(nil, "mods")
		opts := GptOpts{
			Model: stateVar(
				`model`,
				ModelEnv,
				defaultModel,
			),
			NoCache: stateVar(
				`no-cache`,
				NoCacheEnv,
				false,
			),
			Query: strings.Join(args, ` `),
			Quiet: stateVar(`quiet`, QuietEnv, false),
			Role:  stateVar(`role`, RoleEnv, `default`),
			StatusText: stateVar(
				`status-text`,
				StatusTextEnv,
				``,
			),
			Stdin: os.Stdin,
			Title: os.Getenv(TitleEnv),
		}
		return Exec(opts)
	},
}

var commitCmd = &bonzai.Cmd{
	Name:  `commit`,
	Alias: `gc|gptc`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		opts := GptOpts{
			Model: stateVar(
				`model`,
				ModelEnv,
				defaultModel,
			),
			NoCache: true,
			Query:   strings.Join(args, ` `),
			Quiet:   true,
			Role:    `commit-message`,
			StatusText: stateVar(
				`status-text`,
				StatusTextEnv,
				``,
			),
			Stdin: os.Stdin,
			Title: os.Getenv(TitleEnv),
		}
		return Exec(opts)
	},
}

var devCmd = &bonzai.Cmd{
	Name:  `dev`,
	Alias: `code|d`,
	Cmds:  []*bonzai.Cmd{vars.Cmd},
	Do: func(x *bonzai.Cmd, args ...string) error {
		opts := GptOpts{
			Model: stateVar(
				`model`,
				ModelEnv,
				defaultModel,
			),
			NoCache: false,
			Query:   strings.Join(args, ` `),
			Quiet:   false,
			Role:    `dev`,
			StatusText: stateVar(
				`status-text`,
				StatusTextEnv,
				``,
			),
			Stdin: os.Stdin,
			Title: os.Getenv(TitleEnv),
		}
		return Exec(opts)
	},
}

var shellCmd = &bonzai.Cmd{
	Name:  `shell`,
	Alias: `s`,
	Cmds:  []*bonzai.Cmd{vars.Cmd},
	Do: func(x *bonzai.Cmd, args ...string) error {
		opts := GptOpts{
			Model: stateVar(
				`model`,
				ModelEnv,
				defaultModel,
			),
			NoCache: false,
			Query:   strings.Join(args, ` `),
			Quiet:   false,
			Role:    `shell`,
			StatusText: stateVar(
				`status-text`,
				StatusTextEnv,
				``,
			),
			Stdin: os.Stdin,
			Title: os.Getenv(TitleEnv),
		}
		return Exec(opts)
	},
}

var commentCmd = &bonzai.Cmd{
	Name:  `comment`,
	Alias: `c|doc|document`,
	Cmds:  []*bonzai.Cmd{vars.Cmd},
	Do: func(x *bonzai.Cmd, args ...string) error {
		opts := GptOpts{
			Model: stateVar(
				`model`,
				ModelEnv,
				defaultModel,
			),
			NoCache: true,
			Format:  `plain-text`,
			Query:   `create a comment for this function wrapped at 72 and include the function immediately after with no blank line with no markdown or commentary and add square brackets around any symbol reference that could also have a comment except the function name itself and keep it brief`,
			Quiet:   false,
			Role:    ``,
			StatusText: stateVar(
				`status-text`,
				StatusTextEnv,
				``,
			),
			Stdin: os.Stdin,
			Title: os.Getenv(TitleEnv),
		}
		return Exec(opts)
	},
}

var listCmd = &bonzai.Cmd{
	Name:  `list`,
	Alias: `ls`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		convs, err := ListConversations()
		if err != nil {
			return err
		}
		each.Println(convs)
		return nil
	},
}

// stateVar retrieves a value by first checking an environment variable.
// If the environment variable does not exist, it checks bonzai.Vars. If
// neither contain a value, it returns the provided fallback.
func stateVar[T any](key, envVar string, fallback T) T {
	if val, exists := os.LookupEnv(envVar); exists {
		return convertValue(val, fallback)
	}
	if val, err := vars.Data.Get(key); err == nil {
		return convertValue(val, fallback)
	}
	return fallback
}

// convertValue attempts to convert a string to the same type as
// fallback.
func convertValue[T any](val string, fallback T) T {
	var result any = fallback

	switch any(fallback).(type) {
	case string:
		result = val
	case bool:
		result = isTruthy(val)
	case int:
		result, _ = strconv.Atoi(val)
	}

	return result.(T)
}

// isTruthy determines if a string represents a "truthy" value,
// interpreting "t", "true", and positive numbers as true; "f", "false",
// and zero or negative numbers as false.
func isTruthy(val string) bool {
	val = strings.ToLower(strings.TrimSpace(val))
	if slices.Contains([]string{"t", "true"}, val) {
		return true
	}
	if slices.Contains([]string{"f", "false"}, val) {
		return false
	}
	if num, err := strconv.Atoi(val); err == nil {
		return num > 0
	}
	return false
}
