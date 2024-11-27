package tmux

import (
	"context"
	_ "embed"
	"os"

	"github.com/arl/gitmux"
	gtmux "github.com/arl/gitmux/tmux"
	"github.com/arl/gitstatus"
	"github.com/rwxrob/bonzai"
	"gopkg.in/yaml.v3"
)

//go:embed tmux/gitmux.conf
var cfgE []byte

var GitmuxCmd = &bonzai.Cmd{
	Name:    `gitmux`,
	Alias:   `gm`,
	Short:   `tmux plugin for git stats in statusline`,
	MaxArgs: 1,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			args = append(args, ".")
		}
		return Gitmux(args[0])
	},
}

func Gitmux(dir string) (errO error) {
	if dir != "." {
		popdir, err := pushdir(dir)
		if err != nil {
			return err
		}
		defer func() {
			errO = popdir()
		}()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	st, err := gitstatus.NewWithContext(ctx)
	if err != nil {
		return err
	}
	cfg := new(gitmux.Config)
	err = yaml.Unmarshal(cfgE, cfg)
	if err != nil {
		return err
	}
	fmt := &gtmux.Formatter{Config: cfg.Tmux}
	return fmt.Format(os.Stdout, st)
}

func pushdir(dir string) (popdir func() error, err error) {
	pwd := ""
	if pwd, err = os.Getwd(); err != nil {
		return nil, err
	}

	if err = os.Chdir(dir); err != nil {
		return nil, err
	}

	return func() error { return os.Chdir(pwd) }, nil
}
