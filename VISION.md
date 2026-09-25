# Vision: `dot` — A Self-Contained Dotfiles Binary

## The One Idea Everything Else Serves

Config content is embedded into the `dot` binary at build time via
`//go:embed`. A new/target machine needs only that one binary — not a
`git clone`, not a checkout of this repo, not even `git` installed —
to hydrate a full, working dotfiles setup.

The repo itself, the Go toolchain, and GitHub access remain required
for *authoring* configs — that happens on machines the user actively
develops on, where all of that is already available and assumed. `dot`
does not try to eliminate that dependency. What it eliminates is the
need to drag the repo (and git) onto every machine that just wants to
*consume* the configs — a fresh laptop, a short-lived server, a
container.

## Non-Goals

* **Not a package manager.** `dot` delegates actual tool installation
  to the system's package manager (brew, apt, pacman, etc.).
* **Not a secrets manager.** Credentials never live in embedded
  configs.
* **Not fully offline.** Installing real tools still needs network
  access to a package manager. Fetching a *new* binary needs network
  access to GitHub Releases. Neither of these requires `git` or a repo
  checkout, which is the actual guarantee.
* **Not a two-way sync/version-control system.** Earlier drafts of
  this document explored having the binary mutate its own embedded
  config and rebuild itself in place ("commit"), including a
  versioned-history stretch goal. Dropped: it re-derives a chunk of
  what git already does, badly, and the user always has git + the real
  repo on hand anyway. Editing configs means editing the source repo
  and rebuilding — see below.

## Lifecycle

There is exactly one source of truth: the checked-out repo on the
machine you author from. Everything else is downstream of it.

1. **Author.** Edit config source in this repo, on a machine with the
   full toolchain (your dev machine). Commit, push, as normal git
   work.
2. **Build.** On push (to `main`, or on tag — TBD), CI builds the `dot`
   binary per target OS/arch, embedding the current state of the repo,
   and publishes each as a GitHub Release asset.
3. **Deploy to a new machine.** A one-line bootstrap script (`curl |
   sh`) detects OS/arch and downloads the matching release binary. Run
   `dot setup` to hydrate configs onto disk. No git, no clone.
4. **Refresh an existing machine.** `dot upgrade` (naming TBD) fetches
   the latest release binary and replaces the currently running one in
   place — still just an HTTP download of a release asset, never a git
   operation. Followed by `dot setup` to re-hydrate.
5. **Local edits on a deployed machine.** If you hand-edit a live
   config on a machine that only has the binary, that edit is local
   and will not automatically flow back anywhere. If it's worth
   keeping, go make the equivalent edit in the actual repo (step 1) —
   `dot` does not attempt to reconcile this automatically. Drift
   *detection* (below) exists purely so `setup`/`upgrade` don't
   silently destroy such an edit before you've had the chance to
   decide.

## Drift Detection (kept simple)

The one piece of "don't silently overwrite my stuff" behavior worth
keeping, without any of the bidirectional-sync machinery.

### Hydration is a copy, not a symlink

`embed.FS` content lives inside the binary's own memory — there's no
real filesystem path to point a symlink at. So `dot setup` always
writes a real, independent file to disk (`fs.ReadFile` from the
embedded FS → `os.WriteFile` to e.g. `~/.zshrc`). This is consistent
with the existing PRD's move away from fragile symlink-based setups,
not a regression.

### Disambiguating "you edited it" from "upstream changed it"

A hash comparison between the live file and *this binary's* embedded
content is ambiguous across upgrades: a live file can differ from the
current embedded version either because the user hand-edited it, or
because they upgraded to a newer binary whose embedded default simply
changed. Telling these apart requires knowing what was actually
deployed last time, not just what the current binary happens to embed.

So `dot` keeps one small piece of local, per-machine state: a manifest
(e.g. `~/.local/state/dot/manifest.json`, path → hash of what was last
deployed there), written every time `setup` runs. This is deliberately
not the config content itself, and not portable/versioned like the
earlier "history" idea — losing it just degrades to "ask before
overwrite" rather than anything destructive.

* `dot setup <tool>` hashes the live file and compares it against
  **both** the manifest's last-deployed hash and the current embedded
  hash:
  * live == manifest, live != embedded → not drift, just an upstream
    change. Safe to update automatically.
  * live != manifest → user-made drift. **Don't overwrite silently** —
    show a diff and require explicit confirmation (or `--force`) to
    proceed.
  * no manifest entry (first deploy, or file adopted outside `dot`) →
    treat like drift if the file already exists and differs from
    embedded; otherwise just deploy and record it.
* No history, no "commit back into the binary," no per-file version
  log — the manifest holds only a current hash per path, not a log.

### Diffing embedded vs. live content

Both sides end up as plain `[]byte` in memory, so there's no special
embed-vs-disk handling — just two byte slices to compare:

1. Read embedded side: `fs.ReadFile(embeddedFS, "zsh/.zshrc")`.
2. Read live side: `os.ReadFile(filepath.Join(home, ".zshrc"))`.
3. Hash both (sha256) for the cheap equality check used above.
4. When they differ, render a human-readable diff.

For step 4, diff in-process rather than shelling out to system `diff`
(no runtime dependency), rendered as standard unified-diff output.

**Algorithm: histogram diff** (JGit's `HistogramDiff` as reference,
since no existing Go package implements it — all are Myers-family),
falling back to Myers only for unanchorable leftover regions, same as
git/JGit do. Lives in its own standalone Go module, imported by `dot`
like any other dependency.

## Build & Release Pipeline

* CI (`.github/workflows/ci.yml`) builds/vets/tests on every push to
  `main` and every PR, on Linux and macOS.
* Release (`.github/workflows/release.yml`) builds `dot` and `x` for
  linux/darwin × amd64/arm64 and publishes them as GitHub Release
  assets — **tag-triggered only** (`v*.*.*`), not on every push to
  `main`: `dot upgrade` and the bootstrap script both resolve "latest"
  against GitHub Releases, and a release on every commit would make
  "latest" mean "whatever just landed" instead of a deliberate cut.
* Not yet built: version/commit metadata embedded at build time, so
  `dot version` can report exactly which commit's configs it's
  carrying — useful for noticing a deployed machine is stale relative
  to the repo.
* The bootstrap script and `dot upgrade` both resolve "latest" against
  GitHub Releases (a lightweight, unauthenticated HTTP fetch), not the
  git history.

## Open Questions

* Exact naming for `dot upgrade`/`dot setup --force`/diff output
  format.
* Whether `dot setup` needs a `--dry-run` to preview drift across all
  tools at once (a `dot status`-style command), rather than
  discovering drift tool-by-tool.

## Relationship to Existing Docs

* `product_requirements.md` and `docs/prd-cli-restructuring.md` cover
  the broader CLI architecture (command tree, module structure) and
  remain the reference for that.
* `TODO.md`'s "Configuration Tracking Strategy" section sketched an
  early hash-based idea, which this document keeps as the *whole* of
  the tracking story (no history, no bidirectional commit) rather than
  a stepping stone toward something more ambitious.
