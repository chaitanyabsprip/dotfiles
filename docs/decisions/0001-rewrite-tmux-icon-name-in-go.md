# ADR-0001: Rewrite tmux-icon-name as a Go Bonzai command

Date: 2026-03-21
Status: Accepted

## Context

The tmux window auto-rename feature relies on a script (`tmux-icon-name.sh`)
that maps the current pane's process name to an icon and display name. The
script depends on `yq` to parse a YAML config file and is called on every
window rename event (every second by default).

Two problems motivated a rethink:

1. `yq` is slow for this call frequency — it parses the entire YAML file on
   every invocation.
2. The shell script cannot easily express the new matching logic required
   (argument-aware patterns, name substitution, config merging).

The rest of the tmux tooling (sessionizer, session manager, gitmux, etc.) has
already been rewritten as Go Bonzai commands under `dot tmux x`. A new
`IconNameCmd` already exists in `internal/tmux/icon_name.go` but implements
only the original feature set.

### Decision Drivers

- Performance: the command runs on a tight interval and must be fast.
- Consistency: all `dot tmux x` subcommands are Go Bonzai commands.
- Extensibility: new features (pattern matching, name substitution, config
  merge) are awkward to express in shell.
- No new external runtime dependencies.

### Considered Options

#### Option 1: Extend the shell script

- **Pros**: No rewrite, stays simple.
- **Cons**: `yq` parse overhead on every call; complex logic (regex matching,
  config merging) is painful in POSIX shell; inconsistent with the rest of the
  tooling.

#### Option 2: Rewrite in Go as a standalone binary

- **Pros**: Fast, no `yq` dependency, full language expressiveness.
- **Cons**: Separate binary to build and install; not integrated into the
  existing `dot` tree.

#### Option 3: Rewrite in Go as a Bonzai command (chosen)

- **Pros**: Fast, consistent with existing tooling, single binary (`dot`),
  config can be embedded at compile time.
- **Cons**: Slightly more boilerplate than a standalone script.

## Decision

Rewrite `IconNameCmd` in Go as a Bonzai command, replacing both the shell
script (`tmux-icon-name.sh`) and the current partial Go implementation in
`internal/tmux/icon_name.go`. The command lives at `dot tmux x icon`.

## Consequences

### Positive

- No `yq` dependency at runtime; YAML parsed by the Go binary itself.
- Fast startup — compiled binary vs. shell + external process.
- Full Go expressibility for pattern matching and config merge logic.
- Consistent with the rest of `dot tmux x`.

### Negative

- The dotfiles repo must carry the YAML parsing dependency (`gopkg.in/yaml.v3`)
  if not already present.
- The shell script `tmux-icon-name.sh` becomes dead code and should be removed.

## Related Decisions

- ADR-0002: How the full command line is resolved from tmux's `pane_pid`.
- ADR-0003: The YAML config structure that the command parses.
- ADR-0004: How the embedded default config and user config are merged.
