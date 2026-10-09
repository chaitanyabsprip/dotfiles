package env

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rwxrob/bonzai"
	"github.com/rwxrob/bonzai/cmds/help"
	"github.com/rwxrob/bonzai/comp"
	"gopkg.in/yaml.v3"

	"github.com/Chaitanyabsprip/dotfiles/pkg/prompt"
)

var Cmd = &bonzai.Cmd{
	Name:  `env`,
	Short: `work with env files`,
	Long: `
Tools for working with ".env"-style files. See 'env help' for the
commands.`,
	Comp: comp.Cmds,
	Cmds: []*bonzai.Cmd{toJSONCmd, secretProviderCmd, help.Cmd},
	Do: func(x *bonzai.Cmd, _ ...string) error {
		fmt.Printf("%s - %s\n\n", x.Name, x.Short)
		fmt.Println(`COMMANDS:`)
		for _, c := range x.Cmds {
			fmt.Printf("  %-10s - %s\n", c.Name, c.Short)
		}
		return nil
	},
}

var toJSONCmd = &bonzai.Cmd{
	Name:  `tojson`,
	Vers:  `v1.0.0`,
	Short: `convert an env file to json`,
	Usage: `[path]`,
	Long: `
Converts a ".env"-style file to JSON and prints it to stdout. Reads
from stdin when path is "-", or when no path is given and stdin is
piped; otherwise defaults to ".env" in the current directory.`,
	Comp: comp.Opts,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		r, _, closeFn, err := openEnvSource(args, "")
		if err != nil {
			return err
		}
		defer closeFn()

		env, err := parse(r)
		if err != nil {
			return err
		}
		out, err := json.MarshalIndent(env, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
		return nil
	},
}

var secretProviderCmd = &bonzai.Cmd{
	Name:  `secretprovider`,
	Vers:  `v1.0.0`,
	Short: `convert an env file to a SecretProviderClass`,
	Usage: `[path] [name]`,
	Long: `
Converts a ".env"-style file to a Kubernetes SecretProviderClass
manifest (secrets-store.csi.x-k8s.io, AWS Secrets Manager provider)
and prints it to stdout. Each key becomes an object fetched from
Secrets Manager under that same name, synced to a k8s Secret of the
same key. Reads from stdin when path is "-", or when no path is given
and stdin is piped; otherwise defaults to ".env" in the current
directory. name sets metadata.name and defaults to the env file's
basename. Set OUT to a file or directory path to write there instead
of stdout; a directory gets a "secretprovider" file. If OUT's
directory doesn't exist, you'll be prompted to create it.`,
	Comp: comp.Opts,
	Do: func(_ *bonzai.Cmd, args ...string) error {
		var fileArgs []string
		if len(args) > 0 {
			fileArgs = args[:1]
		}
		r, name, closeFn, err := openEnvSource(fileArgs, "")
		if err != nil {
			return err
		}
		defer closeFn()
		if len(args) > 1 {
			name = args[1]
		}

		env, err := parse(r)
		if err != nil {
			return err
		}

		out, err := yaml.Marshal(newSecretProviderClass(name, env))
		if err != nil {
			return err
		}

		w, closeOut, err := openOut(os.Getenv("OUT"), "secretprovider")
		if err != nil {
			return err
		}
		defer closeOut()
		_, err = w.Write(out)
		return err
	},
}

// openOut resolves where output should go. An empty out means stdout. A
// directory (existing, or trailing a path separator) gets defaultName
// appended; anything else is used as the file path as-is. When the target's
// parent directory doesn't exist, the user is prompted to create it.
func openOut(out, defaultName string) (w io.Writer, closeFn func(), _ error) {
	if out == "" {
		return os.Stdout, func() {}, nil
	}

	path := out
	if fi, err := os.Stat(out); (err == nil && fi.IsDir()) ||
		(err != nil && strings.HasSuffix(out, string(filepath.Separator))) {
		path = filepath.Join(out, defaultName)
	}

	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if !prompt.Confirm(os.Stdin, os.Stdout, fmt.Sprintf("%s does not exist, create it?", dir)) {
			return nil, nil, fmt.Errorf("%s does not exist", dir)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, nil, err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}
	return f, func() { f.Close() }, nil
}

type secretProviderClass struct {
	APIVersion string                  `yaml:"apiVersion"`
	Kind       string                  `yaml:"kind"`
	Metadata   map[string]string       `yaml:"metadata"`
	Spec       secretProviderClassSpec `yaml:"spec"`
}

type secretProviderClassSpec struct {
	Provider   string                    `yaml:"provider"`
	Parameters secretProviderClassParams `yaml:"parameters"`
}

type secretProviderClassParams struct {
	Objects string `yaml:"objects"`
}

type secretObject struct {
	ObjectName string `yaml:"objectName"`
	ObjectType string `yaml:"objectType"`
}

func newSecretProviderClass(name string, env map[string]string) secretProviderClass {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	objects := make([]secretObject, 0, len(keys))
	for _, k := range keys {
		objects = append(objects, secretObject{ObjectName: k, ObjectType: "secretsmanager"})
	}
	objectsYAML, _ := yaml.Marshal(objects)

	return secretProviderClass{
		APIVersion: "secrets-store.csi.x-k8s.io/v1",
		Kind:       "SecretProviderClass",
		Metadata:   map[string]string{"name": name},
		Spec: secretProviderClassSpec{
			Provider:   "aws",
			Parameters: secretProviderClassParams{Objects: string(objectsYAML)},
		},
	}
}

// openEnvSource resolves the env file named by args[0] (".env" when args is
// empty), or stdin when that arg is "-" or omitted with stdin piped. name,
// when non-empty, is returned unchanged; otherwise it is derived from the
// file's basename minus extension, or "secrets" for stdin.
func openEnvSource(args []string, name string) (r io.Reader, _ string, closeFn func(), _ error) {
	path := ".env"
	if len(args) > 0 {
		path = args[0]
	}

	if path == "-" || (len(args) == 0 && stdinIsPiped()) {
		if name == "" {
			name = "secrets"
		}
		return os.Stdin, name, func() {}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, "", nil, err
	}
	if name == "" {
		name = strings.TrimPrefix(strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), ".")
		if name == "" {
			name = "secrets"
		}
	}
	return f, name, func() { f.Close() }, nil
}

func stdinIsPiped() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice == 0
}

func parse(r io.Reader) (map[string]string, error) {
	env := map[string]string{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		env[key] = val
	}
	return env, sc.Err()
}
