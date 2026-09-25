# ADR-0012: The x Composing Package Lives at Top-Level `x/`, Not `internal/x/`

Status: Accepted

## Context

Both `docs/prd-cli-restructuring.md` and ADR-0008 describe the `x`
namespace's composing package as `internal/x/x.go`. The actual code has
never lived there — the composing `Cmd` (which pulls together standalone
utilities like `base64`/`depends` and tool `XCmd`s like `tmux.XCmd`) is
`x/x.go`, a top-level, public `package x`, and always has been.

This was an unresolved discrepancy between docs and code, not a decision
anyone had actually made — P2 (extending the `dot x <tool>` pattern to
more tools) needs this settled first, since adding more tools to whichever
location is correct means picking one.

## Decision

The top-level `x/x.go` location is correct. `internal/x/x.go` is not
adopted; the docs that described it were aspirational and are superseded
by this ADR for that one detail (their conceptual content — the
utility/managed-tool split — stands, per ADR-0010).

## Rationale

- **`x` is meant to work as its own standalone binary** (`cmd/x/main.go`),
  per VISION.md's two-binary design — not merely a subcommand tree nested
  inside `dot`. Placing its composing package under `internal/` would
  signal, by Go convention, "not for use outside this module" — the
  opposite of the intent.
- **Consistency with its own sub-packages.** Every utility `x` composes —
  `x/base64`, `x/caseconv`, `x/colors`, `x/depends`, `x/distro`, `x/gpt`,
  `x/have`, `x/install`, `x/workdirs` — already lives at the top level,
  public. Nesting only the composing root under `internal/x/` while every
  child package stays top-level would be an inconsistent split for no
  functional gain (both are importable from anywhere in this module
  either way; it's a single Go module, so `internal/` buys nothing here
  beyond the convention signal above).
- `dot`'s own composing package (`cmd.go`, `package dot` at the repo
  root) follows the same top-level-and-public shape `x/x.go` does — the
  two are peers, not one nested inside the other.

## Consequences

### Positive

- No churn: nothing moves, the PRD/ADR-0008 file paths were simply never
  built and don't need to be.
- P2 (adding more tools' `XCmd`s to the `x` composition) has a settled
  target to extend.

### Negative

- `docs/prd-cli-restructuring.md` and ADR-0008 still show the
  `internal/x/x.go` path in their examples; they're historical records
  and won't be edited retroactively, so a future reader needs this ADR
  to know which path is real. (ADR-0008's `internal/tmux/x.go` example
  is also stale in another way: it shows a nested `RunCmd` grouping that
  never existed in the shipped code, and references a `run.go` that was
  dead code, deleted 2026-09-26.)
