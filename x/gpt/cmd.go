// Package gpt provides a user-friendly CLI for interacting with
// charmbracelet/mods CLI. It is a stateful program such that the users
// will be talking in the same conversation with consecutive calls.
// The package supports different conversation modes, model selection,
// and persistent chat history management.
package gpt

import (
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
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
	Short: `persistent conversation with LLM model using mods`,
	Long: `
Sends args, or stdin when given none, to an LLM via charmbracelet/mods
and prints the reply, continuing the same conversation on later calls.
The model comes from $GPT_MODEL, mods' own default-model config, or
"gemini-free". See 'gpt help' for the specialised roles (commit, dev,
shell, comment) and 'gpt list' for past conversations.`,
	Comp: comp.Combine{comp.Cmds},
	Cmds: []*bonzai.Cmd{
		vars.Cmd,
		commitCmd,
		devCmd,
		shellCmd,
		commentCmd,
		listCmd,
		help.Cmd,
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
	Short: `write a commit message for args or stdin`,
	Long: `
Asks the LLM to write a commit message for args, or stdin when given
none (typically 'git diff --staged'), and prints it. Runs quiet and
uncached: each call is a one-off, not part of the ongoing gpt
conversation.`,
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
	Short: `ask the LLM in its developer role`,
	Long: `
Sends args, or stdin when given none, to the LLM under the "dev" role
(mods role config), for coding questions rather than general chat.`,
	Cmds: []*bonzai.Cmd{vars.Cmd, help.Cmd},
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
	Short: `ask the LLM in its shell role`,
	Long: `
Sends args, or stdin when given none, to the LLM under the "shell"
role (mods role config), for shell command questions.`,
	Cmds: []*bonzai.Cmd{vars.Cmd, help.Cmd},
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
	Short: `write a doc comment for a function`,
	Long: `
Asks the LLM to write a doc comment for the function given in args or
stdin, wrapped at 72 columns, with the function itself immediately
after and any symbol reference bracketed. Runs quiet and uncached.`,
	Cmds: []*bonzai.Cmd{vars.Cmd, help.Cmd},
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
	Short: `list past conversations`,
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
