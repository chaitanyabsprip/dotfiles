# ADR-0004: Embedded default config with user config merge

Date: 2026-03-22
Status: Accepted

## Context

The icon-name command needs a set of default icon mappings (covering common
tools like `nvim`, `git`, `docker`, etc.) that work out of the box without any
user setup. At the same time, users must be able to add their own entries and
patterns, and override the defaults.

Two questions arise:

1. **Where do defaults live?** Options are a file on disk alongside the binary,
   a file in the dotfiles repo that must be symlinked, or a config compiled
   into the binary.
2. **How are user customisations applied?** Options are full replacement of
   defaults, or a merge where user entries win and defaults fill gaps.

### Decision Drivers

- Defaults must work immediately after building the binary, with no extra
  setup steps.
- Users must be able to add new entries and override any default without
  editing the binary's source.
- The defaults should still apply for anything not covered by the user config.

### Considered Options

#### Option 1: Ship defaults as a file on disk only

The binary reads a single config file; if absent, it errors or returns a
fallback icon for everything.

- **Pros**: Simple; user edits one file.
- **Cons**: Requires a setup step (symlinking or copying the file); binary is
  not self-contained; easy to break by deleting the file.

#### Option 2: Embed defaults; user file fully replaces them

`//go:embed icons.yaml` provides defaults. If a user file exists at
`~/.config/tmux/icons.yaml`, it completely replaces the embedded config.

- **Pros**: Binary is self-contained; user has full control.
- **Cons**: User must duplicate all defaults they want to keep; updating
  defaults in the binary has no effect if a user file exists.

#### Option 3: Embed defaults; user file merges with precedence (chosen)

`//go:embed icons.yaml` provides defaults. At startup, both configs are loaded
and merged:

- `config.*` — user values override embedded values key by key.
- `icons` — user entries override embedded entries by key; entries only in the
  embedded config are kept.

If no user file exists, the embedded config is used as-is.

- **Pros**: Binary is self-contained; users only need to write the entries they
  care about; default entries still apply for everything else.
- **Cons**: Merge logic adds complexity. Since priority is determined by longest
  pattern length (see ADR-0003), a user entry with a short pattern may lose to
  an embedded entry with a longer matching pattern.

## Decision

Embed `icons.yaml` into the binary using `//go:embed`. At runtime, merge the
embedded config with `~/.config/tmux/icons.yaml` (if present) using the
following rules:

| Section    | Merge rule |
|------------|------------|
| `config.*` | User value wins per key; missing keys fall back to embedded |
| `icons`    | User entry wins per key; unmatched keys kept from embedded |

## Consequences

### Positive

- Zero-setup: the binary works immediately after installation.
- Users write minimal config — only entries that differ from or extend the defaults.
- Shipping updated defaults in a new binary version automatically benefits
  users who have not overridden those entries.

### Negative

- If a user wants to *remove* a default entry, they cannot do so through the
  merge (only overrides are supported). A workaround (e.g. setting `icon: ""`
  in the user file) may be needed.
- The embedded `icons.yaml` and the user file must stay compatible in schema;
  a schema change in the binary requires user files to be updated.

## Related Decisions

- ADR-0001: Why the command is a Go binary (enables `//go:embed`).
- ADR-0003: The config schema being merged.
