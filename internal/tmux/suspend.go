package tmux

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/run"

	"github.com/Chaitanyabsprip/dotfiles/pkg/tmux"
)

var SuspendCmd = &bonzai.Cmd{
	Name:  `suspend`,
	Alias: `susp`,
	Short: `switch tmux server suspend on or off`,
	Long: `
Binds a key (default F12, or the '@suspend_key' tmux option) to toggle
a "suspended" tmux state: options change, a command runs, and the
prefix key is disabled until the same key resumes it. See 'suspend
help' for the init/on/off commands.`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{
		suspendInitCmd,
		suspendOnCmd,
		suspendOffCmd,
		help.Cmd,
	},
	Do: func(x *bonzai.Cmd, _ ...string) error {
		fmt.Printf("%s - %s\n\n", x.Name, x.Short)
		fmt.Println(`COMMANDS:`)
		for _, c := range x.Cmds {
			fmt.Printf("  %-10s - %s\n", c.Name, c.Short)
		}
		return nil
	},
}

var suspendInitCmd = &bonzai.Cmd{
	Name:  `init`,
	Alias: `setup`,
	Short: `bind the suspend key`,
	Long: `
Binds the suspend key (see 'suspend help') to run 'tmux x suspend on'
and 'tmux x suspend off'. Run this once, e.g. from tmux.conf.`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return SuspendInit()
	},
}

var suspendOnCmd = &bonzai.Cmd{
	Name:    `on`,
	Short:   `suspend the tmux server`,
	Usage:   `<on-suspend-cmd> <suspend-opts>`,
	NumArgs: 2,
	Long: `
Enters the suspended state: disables the prefix key, switches to the
"suspended" key table, cancels any copy mode, un-syncs panes, applies
` + "`<suspend-opts>`" + ` (from '@suspend_suspended_options'), and runs
` + "`<on-suspend-cmd>`" + `. Called by the key binding 'suspend init' sets up.`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return SuspendOn(args[0], args[1])
	},
}

var suspendOffCmd = &bonzai.Cmd{
	Name:    `off`,
	Short:   `resume the tmux server`,
	Usage:   `<on-resume-cmd>`,
	NumArgs: 1,
	Long: `
Reverses 'suspend on': runs ` + "`<on-resume-cmd>`" + `, restores the options
'suspend on' overwrote, restores the prefix key, and clears the key
table. Called by the key binding 'suspend init' sets up.`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		return SuspendOff(args[0])
	},
}

const (
	keyConfig       = `@suspend_key`
	suspendOptsVar  = `@suspend_suspended_options`
	resumeOptsVar   = `@suspend_resumed_options`
	onResumeCmdVar  = `@suspend_on_resume_command`
	onSuspendCmdVar = `@suspend_on_suspend_command`
	defSuspendOpts  = `@mode_indicator_custom_prompt:: ---- , ` +
		`@mode_indicator_custom_mode_style::bg=brightblack,fg=black`
	defOnResumeCmd  = ``
	defOnSuspendCmd = ``
)

func SuspendInit() error {
	key := tmux.GetOption(keyConfig, `F12`)
	suspendOpts := tmux.GetOption(suspendOptsVar, defSuspendOpts)
	onResumeCmd := tmux.GetOption(onResumeCmdVar, defOnResumeCmd)
	onSuspendCmd := tmux.GetOption(onSuspendCmdVar, defOnSuspendCmd)
	exe := run.ExeName()
	run.Out(
		`tmux`, `bind`, `-T`, `root`, key,
		`run-shell`, fmt.Sprintf(
			`%s tmux x suspend on "%s" "%s"`,
			exe, onSuspendCmd, suspendOpts,
		),
	)
	return run.Exec(
		`tmux`, `bind`, `-T`, `suspended`, key,
		`run-shell`, fmt.Sprintf(
			`%s tmux x suspend off "%s"`,
			exe, onResumeCmd,
		),
	)
}

func SuspendOn(onSuspendCmd, suspendOpts string) error {
	prefix := run.Out(`tmux`, `show`, `-qv`, `@suspend_prefix`)
	err := run.Exec(`tmux`, `set`, `-q`, `@suspend_prefix`, prefix)
	if err != nil {
		return err
	}
	err = run.Exec(`tmux`, `set`, `-q`, `prefix`, `none`)
	if err != nil {
		return err
	}
	err = run.Exec(`tmux`, `set`, `key-table`, `suspended`)
	if err != nil {
		return err
	}
	err = run.Exec(
		`tmux`, `if-shell`, `-F`, `#{pane_in_mode}`,
		`send -X cancel`,
	)
	if err != nil {
		return err
	}
	err = run.Exec(
		`tmux`, `if-shell`, `-F`, `#{pane_synchronized}`,
		`"set synchronize-panes off"`,
	)
	if err != nil {
		return err
	}
	setOptionsForSuspendedState(suspendOpts)
	err = run.Exec(`sh`, `-c`, onSuspendCmd)
	if err != nil {
		return err
	}
	return run.Exec(`tmux`, `refresh-client`, `-S`)
}

func setOptionsForSuspendedState(suspendedOpts string) {
	escapedDelim := fmt.Sprintf(
		`%d%d%d`,
		rand.Int(),
		rand.Int(),
		rand.Int(),
	)
	opts := strings.Split(
		strings.ReplaceAll(suspendedOpts, `\\,`, escapedDelim),
		`,`,
	)
	resumedOpts := make([]string, 0)
	for _, opt := range opts {
		opt = strings.TrimSpace(opt)
		if len(opt) == 0 {
			continue
		}
		parts := strings.SplitN(opt, `:`, 3)
		name := strings.TrimSpace(parts[0])
		flags := parts[1]
		value := strings.ReplaceAll(parts[2], escapedDelim, `,`)
		hasValue := len(
			strings.TrimSpace(
				run.Out(
					`tmux`,
					`show`,
					`-qv`+flags,
					name,
				),
			),
		) > 0
		preservedFlags := flags
		if !hasValue {
			preservedFlags += `u`
		}
		preservedValue := strings.TrimSpace(
			run.Out(
				`tmux`,
				`show`,
				`-qv`+flags,
				name,
			),
		)
		opt := fmt.Sprintf(
			`,%s:%s:%s`,
			name,
			preservedFlags,
			strings.ReplaceAll(preservedValue, `,`, `\\,`),
		)
		resumedOpts = append(resumedOpts, opt)
		run.Exec(`tmux`, `set`, `-q`+flags, name, value)
	}
	run.Exec(
		`tmux`, `set`, `-q`, resumeOptsVar,
		strings.Join(resumedOpts, `,`),
	)
}

func SuspendOff(onResumeCommand string) error {
	resumedOpts := run.Out(`tmux`, `show`, `-qv`, resumeOptsVar)
	prefix := strings.TrimSpace(
		run.Out(`tmux`, `show`, `-qv`, `@suspend_prefix`),
	)
	prefixFlags := ``
	if len(prefix) == 0 {
		prefixFlags = `u`
	}
	err := run.Exec(`sh`, `-c`, onResumeCommand)
	if err != nil {
		return err
	}
	setOptionsForResumedState(resumedOpts)
	err = run.Exec(`tmux`, `set`, `-q`+prefixFlags, `prefix`, prefix)
	if err != nil {
		return err
	}
	err = run.Exec(`tmux`, `set`, `-u`, `key-table`)
	if err != nil {
		return err
	}
	return run.Exec(`tmux`, `refresh-client`, `-S`)
}

func setOptionsForResumedState(resumedOptions string) {
	escapedDelim := fmt.Sprintf(
		`%d%d%d`,
		rand.Int(),
		rand.Int(),
		rand.Int(),
	)
	options := strings.Split(
		strings.ReplaceAll(resumedOptions, `\\,`, escapedDelim),
		`,`,
	)
	for _, option := range options {
		option = strings.TrimSpace(option)
		if len(option) == 0 {
			continue
		}
		parts := strings.SplitN(option, `:`, 3)
		name := strings.TrimSpace(parts[0])
		flags := parts[1]
		value := strings.ReplaceAll(parts[2], escapedDelim, `,`)
		run.Exec(`tmux`, `set`, `-q`+flags, name, value)
	}
}
