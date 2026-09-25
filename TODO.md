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

## P1 — Config drift / tracking foundation ✅

Done 2026-09-25, scope = VISION.md's "Drift Detection" section exactly
(hash-based, no history, no bidirectional sync — that section already
settled the "decide scope" question this list used to raise):

- [x] `internal/core/manifest`: hash-based manifest at
      `$XDG_STATE_HOME/dot/manifest.json` (or `~/.local/state/dot/...`),
      path → sha256 of last-deployed content. Tested.
- [x] `internal/core/embed/copy.go`'s `copy()` — the single choke-point
      every tool's `setup` runs through — now does VISION's three-way
      compare (live vs. manifest vs. embedded) instead of `SetupAll`'s old
      behavior of deleting the tool's whole config dir before every
      deploy (which would've destroyed a local edit before drift
      detection ever saw it). Tested (6 cases: first deploy, upstream-only
      change, user drift skipped, forced overwrite, untracked-file
      adoption both ways).
- [x] Diff rendering via `go-udiff` (Myers), already an indirect
      dependency — not the histogram-diff standalone module VISION.md
      describes. `ponytail`: Myers is a materially smaller diff and good
      enough for dotfile-sized configs; upgrade to histogram if diff
      quality on large files actually becomes a problem.
- [x] Force-overwrite via `DOT_FORCE=1` env var — VISION.md left the
      `--force` flag's exact naming as an open question, so this is
      explicitly an interim escape hatch, not the final UX.
- [x] Verified end-to-end in the Debian Docker image: hand-edited
      `~/.config/zsh/.zshrc` survived a second `dot setup zsh` (diff
      printed, file untouched), `DOT_FORCE=1` then overwrote it.
- [x] `dot deps <tool>` implemented for real per ADR-0011 (recursive tree,
      box-drawing chars). `x/depends.Dep`/`PrintTree` + a `toolDeps` map
      in `internal/dot/deps.go`. Populated only for tools whose code
      actually shells out to something extra — **bash → ohmyposh → unzip**,
      **tmux → fzf + tmux itself** — grounded by grepping every
      `exec.Command`/`run.Exec` call in `internal/`, not copied from
      ADR-0011's illustrative example (which shows `ohmyposh` under
      `tmux`; no code evidence supports that edge — tmux doesn't call
      ohmyposh anywhere). A tool with no declared deps renders as a bare
      leaf, not an error — most tools genuinely have none beyond
      themselves. Tested.
- [x] ADR-0011 flipped to `Accepted`, example corrected to match the
      real tree (moved `ohmyposh` from under `tmux` to under `bash`,
      where the code actually puts it).

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
