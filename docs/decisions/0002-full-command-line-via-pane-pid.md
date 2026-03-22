# ADR-0002: Resolve full command line via pane_pid

Date: 2026-03-21
Status: Accepted

## Context

To support argument-aware icon matching (e.g. `nvim -c DBUI` → database icon),
the command needs access to the full command line of the foreground process in
a tmux pane, not just the process name.

tmux exposes two relevant format variables:

- `#{pane_current_command}` — the basename of the foreground process (e.g.
  `nvim`). No arguments.
- `#{pane_pid}` — the PID of the shell that was spawned when the pane was
  created.

`pane_pid` points to the shell, not the foreground process. When a user runs
`nvim -c DBUI`, the shell forks a child; the foreground process is that child.
Running `ps -o args= -p #{pane_pid}` returns the shell invocation (`-zsh`),
not `nvim -c DBUI`.

To get the full command line of the foreground process, the child PID of the
shell must be found first.

Verified on macOS:

```sh
ps -o args= -p $(pgrep -P $(tmux display -t '%10' -p '#{pane_pid}') | tail -1)
# → /usr/local/bin/nvim -c DBUI
```

### Decision Drivers

- Must work on macOS (primary platform).
- Must not require changes to how tmux calls the command (keep `automatic-rename-format` simple).
- Full args must be available, not just the process name.

### Considered Options

#### Option 1: Pass `#{pane_current_command}` only (current behaviour)

- **Pros**: Simple, no process inspection needed.
- **Cons**: No access to arguments; cannot distinguish `nvim` from
  `nvim -c DBUI`.

#### Option 2: Construct the full command line in the tmux format string

Use a shell subcommand inside `automatic-rename-format`:

```tmux
set -g automatic-rename-format \
    "#(ps -o args= -p $(pgrep -P #{pane_pid} | tail -1))"
```

- **Pros**: Binary receives the full string directly.
- **Cons**: Complex, fragile format string; shell interpretation inside tmux
  format strings is error-prone; harder to test.

#### Option 3: Pass `#{pane_pid}` to the binary; resolve internally (chosen)

tmux passes only `pane_pid`. The binary runs:

1. `pgrep -P <pane_pid>` — find child PID of the shell.
2. `ps -o args= -p <child_pid>` — get full command line.

- **Pros**: tmux format string stays simple; all logic is in tested Go code;
  easy to handle edge cases (no child yet, multiple children).
- **Cons**: Two additional process invocations per rename event. Acceptable
  given the rename interval.

## Decision

Pass `#{pane_pid}` from `automatic-rename-format` to `dot tmux x icon`. The
binary resolves the full command line internally via `pgrep` + `ps`.

The normalized command string (basename + args, path stripped) is used for
pattern matching. Example: `/usr/local/bin/nvim -c DBUI` → `nvim -c DBUI`.

## Consequences

### Positive

- `automatic-rename-format` remains a single, readable line.
- All resolution logic is in Go and testable in isolation.
- Path differences across machines (e.g. `/usr/local/bin` vs `/opt/homebrew/bin`)
  are normalised away before matching.

### Negative

- `pgrep` and `ps` must be available on the system (standard on macOS/Linux).
- If a pane has no child process (e.g. shell is idle), the fallback is to use
  `pane_current_command` as the match target.

## Related Decisions

- ADR-0001: Why the command is a Go binary rather than a shell script.
- ADR-0003: How the resolved command string is used for matching.
