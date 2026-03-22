# ADR-0005: Smart defaults config schema with key = pattern

Date: 2026-03-22
Status: Accepted
Supersedes: ADR-0003

## Context

The current config schema has key = display name with optional `pattern` field for
argument-aware matching. This design cannot express "same display name, different
icons" for different commands (e.g., `claude` and `opencode` both showing as
"agent" with different icons).

## Decision Drivers

- Must support same display name with different icons for different commands.
- Must be intuitive to read and write.
- Must be concise for common simple cases.
- Must reduce complexity in matching logic.

## Schema Design

### Syntax

```yaml
config:
  fallback-icon: "?"
  show-name: true

icons:
  # Simple: string value = icon
  # name = key, pattern = exact key
  nvim: "N"
  git: "G"

  # Complex key (spaces or regex chars): object syntax required
  # name = key, pattern = key
  "nvim -c DBUI":
    name: "db"
    icon: "D"

  # Name override: use object syntax
  # pattern = exact key, name = field
  lazygit:
    name: "git"
    icon: "G"

  # Multiple commands, same display name, different icons
  claude:
    name: "agent"
    icon: "C"
  opencode:
    name: "agent"
    icon: "O"
```

### Smart Parsing Rules

| YAML Syntax | Key | Icon | Name | Pattern |
|-------------|-----|------|------|---------|
| `nvim: "N"` | `nvim` | `N` | `nvim` | exact `^nvim$` |
| `lazygit: { name: git, icon: G }` | `lazygit` | `G` | `git` | exact `^lazygit$` |
| `"nvim -c DBUI": { name: db, icon: D }` | `nvim -c DBUI` | `D` | `db` | exact `^nvim -c DBUI$` |
| `"lazygit|gitui": { icon: G }` | `lazygit\|gitui` | `G` | `lazygit\|gitui` | regex `^lazygit\|gitui$` |

### Key Classification

Keys are classified as "simple" or "complex":

- **Simple key**: No spaces, no regex metacharacters (`|`, `.*`, `^`, `$`, `()`, `[]`, `+`, `?`)
  - Pattern = exact match (anchored `^key$`)
  - Name = key (unless `name` field overrides)

- **Complex key**: Has spaces or regex metacharacters
  - Pattern = key as-is (regex match)
  - Name = key (unless `name` field overrides)

### Object Syntax Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | No | Display name. Defaults to key. |
| `icon` | Yes | Icon character. |

## Consequences

### Positive

- Concise for simple entries: `nvim: N`
- Same display name with different icons works cleanly
- Easy to extend embedded entries
- One matching mode — all keys are patterns
- Clear: key = what matches, fields = what displays

### Negative

- Key classification adds magic rules to learn
- Complex keys require object syntax
- Breaking change from current schema

## Related Decisions

- ADR-0003: Original YAML config schema decision (superseded)
- ADR-0004: Embedded default config with user merge
