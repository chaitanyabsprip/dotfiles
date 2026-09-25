# TODO: dot Environment Manager

Reconciles the old `TODO.md` (Phase 0 planning, stale) and `tasks.md`
(P0 task list, more current) into one source of truth. Both originals
are kept at `.archive/{TODO,tasks}.md.bak`.

## Project Status

Verb-first CLI restructuring (`docs/prd-cli-restructuring.md`, ADR-0007-0010)
is implemented and landed: `setup`/`edit`/`install`/`deps`/`status`/`diff`/`init`
compose all 22 tool packages, tiered `slim|quik|full` modes are gone (ADR-0006).

P0 and P1 are both done (below). `internal/dot/install.go` now has an
`InstallCmd` for all 15 tools that need one (bash, bat, ohmyposh, zsh, tmux
plus the 10 added in P0); the other 7 (bin, brew, claude, git, hypr, shell,
vimium) don't install anything — no CLI binary to install, or they're
config-only/self-managing.

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

`tmux`'s `runCmd` question (carried over from old `tasks.md`) — resolved
2026-09-26: `internal/tmux/run.go`'s `runCmd` was a byte-for-byte, unwired
duplicate of `x.go`'s `XCmd` (the one actually composed into `dot x tmux`
via `x/x.go`). Dead code left over from before the verb-first restructure.
Deleted.

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
- [x] `dot status` (list drifted files across all tools) and `dot diff`
      (show their diffs) — added a `CheckMode` to `copy()` so both preview
      drift with zero writes, reusing every tool's existing `SetupCmd`
      instead of new per-tool wiring. Caught and fixed a real bug in the
      process: `bat.SetupCmd` ran `bat cache --build` unconditionally,
      which would've fired as a side effect even under `status`/`diff`'s
      read-only mode.

## P2 — Per-tool `x` utilities beyond tmux ✅

Done 2026-09-26 — turned out to be a two-step check, not a build:

- [x] Target shape confirmed via ADR-0012: top-level `x/x.go` (public
      `package x`) is correct, not the PRD/ADR-0008's `internal/x/x.go` —
      `x` is meant to work as its own standalone binary
      (`cmd/x/main.go`), and every utility sub-package it composes
      (`x/base64`, `x/depends`, etc.) is already top-level and public, so
      nesting only the composing root under `internal/` would be an
      inconsistent split for no functional gain.
- [x] Audited every tool package for non-`Setup`/`Install`/`Edit`/`Deps`
      commands (`grep`, not guessing): `tmux` is genuinely the only one
      with standalone utility scripts today. `bat`'s
      `batGhInstallCmd`/`batPkgInstallCmd` and `shell`'s
      `InstallUnzipCmd` are install-time helpers, not user-facing
      utilities. `claude`'s `HookCmd`/`SessionStartCmd`/`RestoreCmd` are
      already correctly root-level (`dot claude hook ...`, outside
      ADR-0010's managed-tool/utility split entirely) — moving them under
      `x` would break hook paths already configured in people's
      `~/.claude/settings.json`. Nothing to migrate.

## P3 — Polish & release

Lowest priority; do after P0-P2 land, since they change surface area these
would otherwise re-test:

- [ ] `dot upgrade` + GitHub Releases self-update flow — blocked on the
      release pipeline below (needs binaries to actually fetch), and on
      `dot version` (build-time version/commit metadata isn't embedded
      yet, so there's nothing for it to report or compare against).
- [x] CI build/release pipeline. Done 2026-09-26:
      `.github/workflows/ci.yml` (build/vet/test on every push to `main`
      + every PR, Linux and macOS) and `.github/workflows/release.yml`
      (linux/darwin × amd64/arm64 for `dot` and `x`, published as GitHub
      Release assets, tag-triggered on `v*.*.*` — resolves VISION.md's
      "release cadence" open question in favor of tags over push-to-main,
      since `dot upgrade`/the bootstrap script resolve "latest" against
      releases and every-push would make that noisy and undeliberate).
      Verified all 4 target combos actually cross-compile locally before
      committing to them in CI; both workflows pass `actionlint` clean.
- [ ] Static binary build, startup benchmarks
- [ ] Complete test coverage
- [ ] Empty `docs/tasks/` directory — either fill it from this file's
      task breakdown or delete it; don't leave it as dead scaffolding

Explicitly out of scope (PRD "Out of Scope" + VISION.md "Dropped"): tool
groups (`dot setup terminal`), bidirectional `dot pull`/`dot commit` sync.
(`dot status`/`dot diff` were listed here too, but got built as part of
P1 2026-09-26 — read-only drift preview across tools, not the originally-
deferred bidirectional-sync idea, so the PRD's concern didn't actually
apply once scoped that way.)

---

*Reconciled 2026-09-25, superseding TODO.md (2026-03-22) and tasks.md
(2026-03-23).*
