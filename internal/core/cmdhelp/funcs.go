package cmdhelp

import (
	"strings"
	"text/template"

	"github.com/chaitanyabsprip/bonzai"
	"github.com/chaitanyabsprip/bonzai/to"
)

// Funcs is a drop-in replacement for bonzai mark's default FuncMap,
// assignable to a root bonzai.Cmd's Funcs field. It's needed because
// mark depends on the external, unforked github.com/rwxrob/bonzai/mark/funcs
// package, whose functions take *bonzai.Cmd from the upstream module -
// a different type than this fork's *bonzai.Cmd now that the module
// path has been renamed. Overriding every function mark's template
// invokes with one typed against this fork avoids that mismatch.
var Funcs = template.FuncMap{
	"summary":  Summary,
	"usage":    Usage,
	"long":     Long,
	"command":  Command,
	"commands": Commands,
}

// Summary returns the Name joined by a long dash with the Short
// description if it has one.
func Summary(x *bonzai.Cmd) string {
	if x.Short == "" {
		return "`" + x.Name + "`"
	}
	return "`" + x.Name + "` - " + x.Short
}

// Command returns the name of the command joined to any aliases.
func Command(x *bonzai.Cmd) string {
	if x.Alias == "" {
		return x.Name
	}
	return x.Alias + `|` + x.Name
}

// Usage returns Command plus any available Usage information,
// inferring it when absent.
func Usage(x *bonzai.Cmd) string {
	usage := x.Usage
	if usage == "" {
		if len(x.Cmds) > 0 {
			usage = `COMMAND`
		}
		if x.Opts != "" {
			usage += `|` + x.Opts
		}
	}
	return Command(x) + " " + usage
}

// longFuncs excludes "long" itself to avoid an initialization cycle
// with the package-level Funcs map.
var longFuncs = template.FuncMap{
	"summary":  Summary,
	"usage":    Usage,
	"command":  Command,
	"commands": Commands,
}

// Long returns the dedented Long description, filled as a template
// using longFuncs merged with the command's own Funcs.
func Long(x *bonzai.Cmd) (string, error) {
	long := to.Dedented(x.Long)
	f := to.MergedMaps(longFuncs, x.Funcs)
	tmpl, err := template.New("long").Funcs(f).Parse(long)
	if err != nil {
		return "", err
	}
	out := new(strings.Builder)
	if err := tmpl.Execute(out, x); err != nil {
		return "", err
	}
	return out.String(), nil
}
