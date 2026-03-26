# ADR-0008: Module and Namespace Architecture

## Status

Proposed

## Context

The dotfiles system contains two types of functionality:

1. **Tool management**: setup, install, edit, deps for tools (tmux, zsh, git, etc.)
2. **Utility commands**: standalone utilities (base64, colors, sessionizer, etc.)

These need to coexist cleanly in the CLI while maintaining code
organization and discoverability.

## Decision

### Tool Modules as Verb Command Sources

Each tool package in `internal/<tool>` exports verb-specific commands
(`SetupCmd`, `InstallCmd`, `EditCmd`, `DepsCmd`) that are composed into
the corresponding top-level verb commands.

**Composition in internal/dot/setup.go:**

```go
var SetupCmd = &bonzai.Cmd{
    Name: `setup`,
    Cmds: []*bonzai.Cmd{
        tmux.SetupCmd,
        zsh.SetupCmd,
        git.SetupCmd,
        // ...
    },
}
```

### XCmd: Tool Utility Subtree

Each tool package exports an `XCmd` — a self-contained subtree of its
utility commands. This keeps tool packages as the single owner of their
utility subtree.

**Example - internal/tmux/x.go:**

```go
var XCmd = &bonzai.Cmd{
    Name: `tmux`,
    Cmds: []*bonzai.Cmd{
        SessionizerCmd,
        IconCmd,
        RunCmd,
    },
}
```

### X Namespace Composition

The `x` namespace imports and composes `XCmd` from each tool package:

**internal/x/x.go:**

```go
var Cmd = &bonzai.Cmd{
    Name: `x`,
    Cmds: []*bonzai.Cmd{
        tmux.XCmd,
        zsh.XCmd,
        base64.Cmd,
        colors.Cmd,
        creashCmd,
        catcCmd,
        depends.Cmd,
        distro.Cmd,
        gpt.Cmd,
        have.Cmd,
        workdirs.Cmd,
        caseconv.Cmd,
        help.Cmd,
    },
}
```

### Command Structure

```bash
dot x tmux sessionizer    # SessionizerCmd from internal/tmux/sessionizer.go
dot x tmux icon          # IconCmd from internal/tmux/icon.go
dot x tmux run           # RunCmd from internal/tmux/run.go
dot x zsh <cmd>          # Future zsh utilities
```

### Remove x install Subcommand

The `x install` subcommand is removed. Tool installation logic resides
in `internal/<tool>/InstallCmd`. Users install tools via `dot install
<tool>`.

### X Namespace Utilities

The `x` namespace contains standalone utility commands that don't belong
to a specific tool:

| Command | Description |
| --------- | ------------- |
| `x base64` | Base64 encode/decode |
| `x colors` | Color utilities |
| `x creash` | Create shell scripts |
| `x catc` | Cat script contents |
| `x depends` | Check dependencies |
| `x distro` | OS detection |
| `x gpt` | LLM integration |
| `x have` | Check executables |
| `x workdirs` | Project directory management |
| `x caseconv` | Case conversion |

## Consequences

### Positive

- **Clear ownership**: Tool packages own their setup/install/edit/deps logic
- **X namespace organization**: Tool-specific utilities grouped by tool
  (`x tmux <cmd>`)
- **Consistent tool exposure**: All tools follow the same pattern
- **No duplication**: Utilities live in one place, exposed via XCmd
- **Adding utilities is local**: New tmux utility only touches `internal/tmux/`

### Negative

- **Import complexity**: XCmd must import all tool utility packages

## Notes

- All tool-specific utilities must be exposed under `x` (none hidden)
- `x` exists as both standalone binary (`cmd/x/main.go`) and `dot x` subcommand
