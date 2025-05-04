package tmux

import (
	"fmt"
	"path/filepath"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/is"
	"github.com/rwxrob/bonzai/yq"

	"github.com/Chaitanyabsprip/dotfiles/pkg/env"
)

var IconNameCmd = &bonzai.Cmd{
	Name:  `icon`,
	Alias: `iname|i|name`,
	Short: `display appropriate name with icon for process`,
	Do: func(x *bonzai.Cmd, args ...string) error {
		if len(args) == 0 {
			args = append(args, ``)
		}
		return IconName(args[0])
	},
}

func IconName(name string) error {
	cfgFile := filepath.Join(
		env.HOME,
		`.config`,
		`tmux`,
		`icons.yml`,
	)
	icon, err := yq.EvaluateToString(
		fmt.Sprintf(`.icons.%s`, name),
		cfgFile,
	)
	if err != nil {
		return err
	}
	if len(icon) == 0 {
		icon, err = yq.EvaluateToString(
			`.config.fallback-icon`,
			cfgFile,
		)
		if err != nil {
			return err
		}
	}
	icon = icon[1 : len(icon)-1]
	showName, err := yq.EvaluateToString(
		`.config.show-name`,
		cfgFile,
	)
	if err != nil {
		return err
	}
	if is.Truthy(showName) {
		fmt.Println(icon, name)
	} else {
		fmt.Println(icon)
	}
	return nil
}
