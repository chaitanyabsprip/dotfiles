# ADR-0007: CLI Command Tree Structure

## Status

Proposed

## Context

The current CLI structure uses noun-first organization where tools are top-level commands (`dot tmux`, `dot zsh`, `dot brew`) with verbs as subcommands (`dot tmux init`, `dot zsh setup`). This was established before fully considering the trade-offs between discoverability, tab-completion ergonomics, and cross-tool operations.

## Decision

We will adopt a **verb-first** command tree structure at the top level.

### Top-Level Commands

| Command | Type | Description |
|---------|------|-------------|
| `setup` | Verb | Hydrates configuration files |
| `install` | Verb | Installs external dependencies |
| `edit` | Verb | Opens tool configuration in editor |
| `deps` | Verb | Displays dependency tree for tool |
| `init` | Meta-command | Full bootstrap: setup + install for all |
| `x` | Namespace | Utility commands |

### Bare Verb Behavior

| Command | No Arg | With Tool |
|---------|--------|-----------|
| `dot setup` | Runs all tools | Runs specific tool |
| `dot install` | **Errors** | Runs specific tool |
| `dot edit` | Errors | Edits specific tool |
| `dot deps` | Errors | Shows specific tool |
| `dot init` | Full bootstrap | N/A |

### Command Examples

```bash
# Specific tool operations
dot setup tmux          # Hydrate tmux config
dot install tmux        # Install tmux binary
dot edit tmux           # Edit tmux config
dot deps tmux           # Show tmux dependencies

# All tools
dot setup              # Hydrate all configs (safe, idempotent)
dot install all        # Install all binaries (explicit intent)
dot init               # Full bootstrap

# Utilities (see ADR-0008)
dot x tmux sessionizer
dot x base64 encode
```

### File Structure

```
cmd/
  dot/main.go           # Entry point: dot.Cmd.Exec()
  x/main.go             # Entry point: x.Cmd.Exec()

internal/
  dot/
    setup.go            # SetupCmd
    install.go          # InstallCmd
    edit.go             # EditCmd
    deps.go             # DepsCmd
    init.go             # InitCmd

  x/
    x.go                # XCmd

  tmux/
    x.go                # XCmd (subtree for utilities)
    setup.go            # SetupCmd
    install.go          # InstallCmd
    edit.go             # EditCmd
    deps.go             # DepsCmd
    sessionizer.go      # SessionizerCmd
    icon.go             # IconCmd
    run.go              # RunCmd

  zsh/
    x.go                # XCmd
    setup.go
    install.go
    ...

  (other tools)
```

## Consequences

### Positive

- **Discoverability**: Tab-completing `dot setup<TAB>` shows all available tools
- **Consistency**: Single pattern for all tool operations
- **Scriptability**: `dot setup tmux` reads clearly in scripts

### Negative

- **Breaking change**: Existing `dot tmux init` commands must be updated to `dot setup tmux`

## Notes

- `status` and `diff` verbs are deferred to future implementation
- Tool utility subcommands under `x` follow different conventions (see ADR-0008)
