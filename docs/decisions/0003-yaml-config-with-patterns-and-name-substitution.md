# ADR-0003: YAML config with pattern matching and name substitution

Date: 2026-03-22
Status: Superseded

## Superseded by

[ADR-0005](0005-flip-config-schema.md): Change config schema — key = pattern, fields = name + icon

## Context

The original `icon_names.yml` is a flat map of process name → icon string.
This works for simple cases but has two gaps:

1. **No argument-aware matching.** `nvim` and `nvim -c DBUI` are
   indistinguishable. A shell alias `db=vim -c DBUI` was attempted as a
   workaround but aliases are expanded before tmux inspects the process, so
   the alias name is lost by the time the pane command is read.

2. **No name substitution.** The display name is always the raw process name.
   There is no way to show `db` instead of `nvim` when DBUI is running, or
   `git` instead of `lazygit`.

The config format must be extended to express both capabilities while remaining
hand-editable YAML.

### Decision Drivers

- Must support regex matching on the full command line (basename + args).
- Must support simple exact matching on process name for the common case.
- The display name must be decoupled from the process name.
- More specific matches must win when multiple entries match the same command.
- Must stay readable without tooling.

### Considered Options

#### Option 1: Flat map with argument suffix keys

Extend the existing map with keys like `"nvim -c DBUI"`:

```yaml
icons:
  nvim: ""
  "nvim -c DBUI": ""
```

- **Pros**: Minimal format change.
- **Cons**: No way to express custom display names; no regex support; key
  order in YAML maps is undefined so priority is ambiguous.

#### Option 2: Separate `patterns` list alongside the `icons` map

Keep a simple `icons` map for process defaults and add an ordered `patterns`
list for argument-aware rules:

```yaml
icons:
  nvim:
    icon: ""
    name: nvim

patterns:
  - match: "nvim -c DBUI"
    icon: ""
    name: db
```

- **Pros**: Clear separation of concerns; patterns are explicitly ordered.
- **Cons**: Two sections to maintain; the `name` field is redundant on
  `icons` entries where the key already names the process; users must
  understand the two-section mental model.

#### Option 3: Single unified `icons` map with optional `pattern` field (chosen)

All entries live in one `icons` map. The key is always the display name. Each
entry has an optional `icon` and an optional `pattern`; at least one must be set.
Set `icon` to `""` to display no icon:

```yaml
config:
  fallback-icon: "?"
  show-name: true

icons:
  nvim:
    icon: ""
  db:
    icon: ""
    pattern: "nvim -c DBUI"
  git:
    icon: ""
    pattern: "lazygit|gitui"
```

- Entries **without** `pattern`: key matched exactly against
  `#{pane_current_command}`.
- Entries **with** `pattern`: regex matched against the normalized full
  command line (`<basename> <args>`).
- When multiple entries match, the one with the longest `pattern` string wins
  (most specific match).
- **Pros**: One section, one mental model; key is unambiguously the display
  name; no separate `name` field needed; specificity-based priority removes
  the need for explicit ordering.
- **Cons**: YAML maps have no guaranteed order, but since priority is
  determined by pattern length rather than position, order does not matter.

## Decision

Use a single `icons` map where the key is the display name. Both `icon`
and `pattern` are optional, but at least one must be set. Set `icon` to `""`
for no icon. Entries without a `pattern` match via `#{pane_current_command}`;
entries with a `pattern` match via the normalized full command line. The
longest matching pattern wins.

## Consequences

### Positive

- One section to learn and maintain.
- The key is always the display name — no redundant `name` field.
- Regex patterns cover all argument-aware cases.
- Specificity-based priority (longest pattern wins) is intuitive and requires
  no explicit ordering by the user.

### Negative

- The `icons` map values change from plain strings to objects — a breaking
  change from the original `icon_names.yml` format.
- Entries without a `pattern` and entries with a `pattern` use different match
  targets (`pane_current_command` vs. full command line), which must be clearly
  documented to avoid confusion.

## Related Decisions

- ADR-0002: How the normalized command string used for pattern matching is
  produced.
- ADR-0004: How user-defined entries are merged with the embedded defaults.
