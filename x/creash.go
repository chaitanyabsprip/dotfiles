package x

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rwxrob/bonzai/comp"
	"github.com/rwxrob/bonzai/futil"
	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/run"
)

var creashCmd = &bonzai.Cmd{
	Name:    `creash`,
	Alias:   "csh",
	Vers:    "v1.0.0",
	Short:   `Create a new shell script in pwd`,
	Long:    `Create a new shell script in pwd`,
	MinArgs: 1,
	Comp:    comp.Opts,
	Call: func(x *bonzai.Cmd, args ...string) error {
		edit := len(os.Getenv("EDIT")) > 0
		creash(edit, args...)
		return nil
	},
}

func creash(edit bool, names ...string) {
	scriptsDir, err := os.Getwd()
	if err != nil {
		return
	}
	for _, name := range names {
		path := filepath.Join(scriptsDir, name)
		if futil.Exists(path) {
			fmt.Print("File already exists: ", name)
			continue
		}
		template := "#!/bin/sh\n"
		if err := os.WriteFile(path, []byte(template), 0o755); err != nil {
			fmt.Print("Could not create file: ", name)
			continue
		}
	}
	if edit {
		run.Exec(append([]string{editor()}, names...)...)
	}
}

func editor() string {
	ed := os.Getenv("VISUAL")
	if len(ed) == 0 {
		ed = os.Getenv("EDITOR")
	}
	if len(ed) == 0 {
		ed = "nvim"
	}
	if len(ed) == 0 {
		ed = "vim"
	}
	if len(ed) == 0 {
		ed = "vi"
	}
	return ed
}

func scriptsDir() string {
	scriptsDir := os.Getenv("SCRIPTS")
	if len(scriptsDir) == 0 {
		scriptsDir = filepath.Join(
			os.Getenv("HOME"),
			".local",
			"bin",
		)
	}
	return scriptsDir
}
