# ADR-0010: X Utilities Are Not Managed Tools

Status: Proposed

## Context

The `x` namespace contains utility commands (`base64`, `colors`, `caseconv`, `depends`, `distro`, `gpt`, `have`, `workdirs`, `creash`, `catc`) alongside tool-specific utilities (`tmux sessionizer`, `tmux icon`, `tmux run`). These utilities serve a fundamentally different purpose from managed tools like `tmux`, `zsh`, or `git`.

With the verb-first command tree (ADR-0007), it needs to be clear which commands participate in the top-level verbs (`setup`, `install`, `edit`, `deps`) and which do not.

## Decision

Commands under the `x` namespace are **utility commands**, not managed tools. They do not participate in the `setup`, `install`, `edit`, or `deps` verbs.

### Distinction

| Concern | Managed Tools (`internal/<tool>`) | X Utilities (`x/`, `internal/<tool>/XCmd`) |
|---------|-----------------------------------|---------------------------------------------|
| Example | tmux, zsh, git, brew | base64, colors, caseconv, tmux sessionizer |
| Setup | `dot setup tmux` | N/A |
| Install | `dot install tmux` | N/A |
| Edit | `dot edit tmux` | N/A |
| Deps | `dot deps tmux` | N/A |
| Invocation | Via top-level verbs | `dot x <cmd>` or `x <cmd>` |

### Rationale

- Utilities are self-contained commands with no configuration files to hydrate or external dependencies to install.
- Including them in verb commands would pollute tab-completion (e.g., `dot setup base64` makes no sense).
- The `x` namespace is a collection of developer productivity tools, not environment configuration.

## Consequences

### Positive

- Clean separation between environment management and utility commands
- Top-level verbs only show meaningful targets
- No confusion about what `dot setup` operates on

### Negative

- If a utility ever needs setup/install semantics, it must be promoted to a managed tool in `internal/<tool>`
