# ADR-0009: Init Command Behavior

## Status

Proposed

## Supersedes

ADR-0006: Full Setup Only - No Tiered Configuration

## Context

ADR-0006 established that dotfiles provisioning should be "full setup only" without tiered configuration options (slim, quik, full). The implementation of `dot init` needs to be clarified with the new verb-first command structure.

## Previous Decision (ADR-0006)

> The dotfiles provisioning should be "full setup only" without tiered configuration options. Users should not be presented with choices between different levels of configuration at setup time.

## Updated Decision

The `dot init` command performs a full environment bootstrap by executing:
1. `setup` for all configured tools (hydrates all configuration files)
2. `install` for all configured tools (installs all external dependencies)

```bash
# Equivalent to:
dot setup && dot install all
```

### Use Cases

| Scenario | Command |
|----------|---------|
| New machine setup | `dot init` |
| SSH into server, config only | `dot setup tmux` |
| SSH into server, install binary | `dot install tmux` |
| Update all configs | `dot setup` |
| Reinstall all binaries | `dot install all` |

### Verb Behavior Summary

| Command | No Arg | Tool Arg |
|---------|--------|----------|
| `dot setup` | Runs all tools | Runs specific tool |
| `dot install` | Errors | Runs specific tool |
| `dot edit` | Errors | Edits specific tool |
| `dot deps` | Errors | Shows specific tool |
| `dot init` | Full bootstrap | N/A |

**Rationale:**
- `dot setup` with no arg runs all tools because config hydration is safe, idempotent, and is the common case
- `dot install` with no arg errors because installing all binaries is heavyweight and requires explicit intent (`all`)
- `dot edit` and `dot deps` require a tool arg because operating on "all" doesn't make semantic sense

## Consequences

### Positive

- **Single command bootstrap**: `dot init` provisions entire environment
- **Clear mental model**: `setup` = configs, `install` = binaries
- **Explicit intent for install**: Users must opt-in to full binary installation

### Negative

- **All-or-nothing**: Users cannot selectively install tools during `dot init`

## Notes

- This ADR refines ADR-0006 to explicitly include the `install` phase alongside `setup`
- Future work may explore tool groups (e.g., `dot setup terminal` for terminal-related tools only)
