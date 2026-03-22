# ADR-0006: Full Setup Only

Date: 2026-03-23
Status: Accepted

## Context

The original PRD specified tiered configuration modes (slim|quik|full)
for setup commands, allowing users to choose different levels of
configuration depth. This was implemented across all modules with the
`Opts: slim|quik|full` option.

### Rationale

1. **Simplicity**: Single setup mode reduces cognitive overhead for users
2. **Default to complete**: Users typically want full configuration when
   running setup
3. **Less maintenance**: No need to maintain multiple code paths for different levels
4. **Deferred granularity**: Add tiered options only when specifically
   requested by users

### Changes Required

1. Remove `Opts: slim|quik|full` from all module setup commands
2. Remove mode selection logic from zsh initCmd and tmux initCmd
3. Simplify setup commands to always perform full setup
4. Update tasks.md to reflect simplified scope

## Decision

We will implement only "full" setup for each tool. The setup commands
will perform complete configuration deployment without any mode
selection options.

## Consequences

### Positive

- Simplified command interface
- Reduced code complexity
- Easier maintenance
- Clearer user experience

### Negative

- Users cannot opt for minimal setup if desired in the future
- Need to implement separate PR/ADR to add granularity later

## Implementation

Completed: 2026-03-23

- Removed `Opts` and `Comp: comp.Opts` from 21 module setup commands
- Simplified tmux initCmd to always run full setup
- Removed debug print statements from hypr module

## Related Decisions

- ADR-0004: Embedded default config with user config merge (foundational)
- Future ADR: Add tiered configuration (if needed)
