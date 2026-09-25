# TODO: dot Environment Manager

Reconciles the old `TODO.md` (Phase 0 planning, stale) and `tasks.md`
(P0 task list, more current) into one source of truth. Both originals
are kept at `.archive/{TODO,tasks}.md.bak`.

## Project Status

Verb-first CLI restructuring (`docs/prd-cli-restructuring.md`, ADR-0007-0010)
is implemented and landed: `setup`/`edit`/`install`/`deps`/`init` compose all
22 tool packages, tiered `slim|quik|full` modes are gone (ADR-0006).

Confirmed against `internal/dot/install.go` on 2026-09-25:
tool packages with a working `InstallCmd`: **bash, bat, ohmyposh, zsh, tmux**.

---

## P0 — Finish install command coverage ✅

Done 2026-09-25. Added `InstallCmd` for each, wired into
`internal/dot/install.go`'s `InstallCmds` slice, and validated via Docker
(`Dockerfile.arch-e2e` for pacman, `Dockerfile` for apt):

- [x] alacritty *(brew: cask)*
- [x] fish
- [x] gh *(pacman: `github-cli`; apt/dnf/brew: `gh`)*
- [x] gitui *(unavailable via apt — prints manual-install hint, non-fatal)*
- [x] kitty *(brew: cask)*
- [x] lsd
- [x] sqlfluff *(unavailable via dnf — same hint path)*
- [x] starship *(unavailable via apt/dnf — uses official curl installer,
      same pattern as `ohmyposh`)*
- [x] waybar *(Linux-only, no brew formula)*
- [x] dirs *(package is `xdg-user-dirs`, Linux-only)*

Added `x/install.Pkg()` to share the pacman/apt/dnf/brew switch across
these 8 straightforward tools instead of copy-pasting it 10 times.

No install needed (skip, per `x/install` + PRD): bash ✅, git, bin, shell,
vimium.

Also resolve: does `tmux`'s `runCmd` still belong now that setup is
full-only, or should it fold into `InstallCmd`/`SetupCmd`? (Carried over
from old `tasks.md`, still open.)

## P1 — Config drift / tracking foundation

The single biggest gap between VISION.md and the codebase — no `manifest`,
`hydrate`, or `drift` code exists anywhere yet, and everything downstream
(`dot deps`, `dot pull`/`diff`) depends on deciding this first.

- [ ] Decide scope: hash-based MVP (SHA256 per deployed file, recorded at
      `dot setup` time) vs. the fuller state-machine from `thoughts.md`/
      `notes.dump.md` (both archived — those proposed `dot pull`/`dot commit`
      bidirectional sync, which VISION.md marks **dropped** as a non-goal;
      don't resurrect that scope without deciding it's back in)
- [ ] Land state storage location (`~/.local/state/dot/dot.lock` was the
      prior working assumption)
- [ ] Implement `dot deps <tool>` for real — currently
      `internal/dot/deps.go` is a hardcoded stub ("no dependencies
      defined"). ADR-0011 (recursive dependency tree view) is written but
      **Proposed**, not implemented — implement it or downgrade the ADR.

## P2 — Per-tool `x` utilities beyond tmux

PRD user stories 7/9/15 want every tool's utilities reachable as
`dot x <tool> <cmd>`, organized by tool prefix. Only `internal/tmux/x.go`
does this today.

- [ ] Confirm target shape: PRD's proposed `internal/x/x.go` composing
      package doesn't exist — top-level `x/x.go` currently plays that role
      with a different layout. Pick one and document it (update ADR-0010
      or add a new one) before adding more tools here.
- [ ] Migrate/add tool-specific utilities for tools that have them today
      outside the `x` namespace (audit `internal/<tool>/` for anything
      script-like first — don't invent utilities that don't exist yet).

## P3 — Polish & release

Lowest priority; do after P0-P2 land, since they change surface area these
would otherwise re-test:

- [ ] `dot upgrade` + GitHub Releases self-update flow
- [ ] CI build/release pipeline (per-OS/arch binaries on push/tag)
- [ ] Static binary build, startup benchmarks
- [ ] Complete test coverage
- [ ] Empty `docs/tasks/` directory — either fill it from this file's
      task breakdown or delete it; don't leave it as dead scaffolding

Explicitly out of scope (PRD "Out of Scope" + VISION.md "Dropped"):
`dot status`/`dot diff` as originally conceived, tool groups
(`dot setup terminal`), bidirectional `dot pull`/`dot commit` sync.

---

*Reconciled 2026-09-25, superseding TODO.md (2026-03-22) and tasks.md
(2026-03-23).*
