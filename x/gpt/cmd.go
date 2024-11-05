package gpt

import (
	"os"
	"strings"

	"github.com/Chaitanyabsprip/dot/x/depends"
	"github.com/rwxrob/bonzai"
)

var Cmd = &bonzai.Cmd{
	Name: `gpt`,
	Init: func(x *bonzai.Cmd, args ...string) error {
		depends.On(nil, "mods")
		return nil
	},
	// Cmds: []*bonzai.Cmd{vars.Cmd},
	Call: func(x *bonzai.Cmd, args ...string) error {
		opts := GptOpts{
			Model:      os.Getenv(`GPT_MODEL`),
			NoCache:    false,
			Query:      strings.Join(args, ` `),
			Quiet:      false,
			Role:       os.Getenv(`GPT_ROLE`),
			StatusText: `Ummm`,
			Stdin:      os.Stdin,
			Title:      ``,
		}
		return Exec(opts)
	},
}

var CommitCmd = &bonzai.Cmd{
	Name: `gptc`,
	Init: func(x *bonzai.Cmd, args ...string) error {
		depends.On(nil, "mods")
		return nil
	},
	// Cmds: []*bonzai.Cmd{vars.Cmd},
	Call: func(x *bonzai.Cmd, args ...string) error {
		opts := GptOpts{
			Model:      os.Getenv(`GPT_MODEL`),
			NoCache:    true,
			Query:      strings.Join(args, ` `),
			Quiet:      true,
			Role:       `commit-message`,
			StatusText: `Ummm`,
			Stdin:      os.Stdin,
			Title:      ``,
		}
		return Exec(opts)
	},
}

var DevCmd = &bonzai.Cmd{
	Name: `gptd`,
	Init: func(x *bonzai.Cmd, args ...string) error {
		depends.On(nil, "mods")
		return nil
	},
	// Cmds: []*bonzai.Cmd{vars.Cmd},
	Call: func(x *bonzai.Cmd, args ...string) error {
		opts := GptOpts{
			Model:      os.Getenv(`GPT_MODEL`),
			NoCache:    false,
			Query:      strings.Join(args, ` `),
			Quiet:      false,
			Role:       `dev`,
			StatusText: `Ummm`,
			Stdin:      os.Stdin,
			Title:      ``,
		}
		return Exec(opts)
	},
}

var ShellCmd = &bonzai.Cmd{
	Name: `gpts`,
	Init: func(x *bonzai.Cmd, args ...string) error {
		depends.On(nil, "mods")
		return nil
	},
	// Cmds: []*bonzai.Cmd{vars.Cmd},
	Call: func(x *bonzai.Cmd, args ...string) error {
		opts := GptOpts{
			Model:      os.Getenv(`GPT_MODEL`),
			Query:      strings.Join(args, ` `),
			Role:       `shell`,
			StatusText: `Ummm`,
			Title:      ``,
			NoCache:    false,
			Quiet:      false,
			Stdin:      os.Stdin,
		}
		return Exec(opts)
	},
}

var CommentCmd = &bonzai.Cmd{
	Name: `gpte`,
	Init: func(x *bonzai.Cmd, args ...string) error {
		depends.On(nil, "mods")
		return nil
	},
	// Cmds: []*bonzai.Cmd{vars.Cmd},
	Call: func(x *bonzai.Cmd, args ...string) error {
		opts := GptOpts{
			Model:      os.Getenv(`GPT_MODEL`),
			NoCache:    false,
			Format:     `comment`,
			Query:      strings.Join(args, ` `),
			Quiet:      false,
			Role:       os.Getenv(`GPT_ROLE`),
			StatusText: `Ummm`,
			Stdin:      os.Stdin,
			Title:      ``,
		}
		return Exec(opts)
	},
}
