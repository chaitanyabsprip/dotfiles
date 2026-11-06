package dot

import (
	"testing"

	"github.com/chaitanyabsprip/bonzai"
)

// TestAliasesUnique guards against two sibling commands answering to
// the same name or alias, where bonzai would silently pick one.
func TestAliasesUnique(t *testing.T) {
	var walk func(parent *bonzai.Cmd, path string)
	walk = func(parent *bonzai.Cmd, path string) {
		seen := map[string]string{}
		for _, c := range parent.Cmds {
			for _, n := range c.Aliases() {
				if prev, ok := seen[n]; ok && prev != c.Name {
					t.Errorf("%s: %q is used by both %s and %s", path, n, prev, c.Name)
				}
				seen[n] = c.Name
			}
			walk(c, path+" "+c.Name)
		}
	}
	walk(Cmd, Cmd.Name)
}
