# Dotfiles

Terminal-centric config for macOS and Linux: zsh/bash/fish sharing one core,
tmux, Neovim, Ghostty, git, and ~150 scripts in `bin/`.

## Bootstrap

```sh
git clone --recurse-submodules git@github.com:Chaitanyabsprip/dotfiles.git ~/dotfiles
~/dotfiles/dotme serve   # symlink every config listed in .dotme.lock
~/dotfiles/setup         # install tools: Brewfile on macOS, package manager on Linux
```

The repo must live at `~/dotfiles`; the shell configs and a few others
assume that path. Both commands are safe to re-run: `dotme serve` skips links
that are already correct and moves anything else at the destination to
`<path>.bak.<timestamp>`; `setup` skips installed tools. `./setup <target>`
runs one part (`packages`, `zsh`, `tmux`, `nvim`, `mise`, `go`, ...).

## Layout

| Path | What |
|---|---|
| `shell/` | Shared by zsh, bash and fish: `env.sh` (environment, PATH), `aliases`, `functions.sh`, `fzf.sh` |
| `zsh/`, `bash/`, `fish/` | Per-shell entry points; only shell builtins and plugins live here |
| `bin/` | Scripts, on PATH via `shell/env.sh`; indexed in [`bin/README.md`](bin/README.md) |
| `tmux/` | `tmux.conf` plus `conf.d/` (statusline, harpoon, resurrect, ...) |
| `nvim/` | Neovim config (git submodule) |
| `git/`, `gh/`, `ghostty/`, `mise/`, `bat/`, ... | One directory per tool, linked into `~/.config` |
| `hypr/`, `waybar/`, `dirs` | Linux desktop only (`: linux` in `.dotme.lock`) |
| `Brewfile` | macOS packages (`brew bundle dump --force --no-vscode --no-go --no-cargo --no-uv --no-npm`) |
| `dotme`, `.dotme.lock` | The symlink manager and its `source: destination[: os]` list |
| `dockerfiles/` | Linux environments that bootstrap with `dotme serve && setup` |

## Conventions

- **One shell core.** `shell/env.sh` and `shell/aliases` are written in a
  POSIX subset that fish can source too (rules at the top of `env.sh`).
  Add an alias once and every shell gets it. Anything that doesn't need to
  change the current shell goes in `bin/` as a script instead of a function.
- **Runtimes come from mise** (`mise/config.toml`): `mise use -g node@22`.
- **Local, untracked overrides:** `zsh/conf.d/work.zsh`,
  `bash/conf.d/local.bash`, `fish/conf.d/local.fish`, `git/local` (e.g. a
  work email), `~/.bash_profile.local`.
- **`~/.local/bin` is not the repo.** It holds installer drops and scratch
  scripts; `$DOTFILES/bin` comes first on PATH.

## Trying new things

- **Shell:** drop a file in `zsh/conf.d/` (or `bash/conf.d/`,
  `fish/conf.d/`); rename it to `*.off` to disable it. For a throwaway zsh
  setup: `ZDOTDIR=/tmp/try zsh`.
- **Scripts:** `hem -f <name> [command]` creates one in `~/.local/bin`; move
  it to `bin/` once it earns its place.
- **A new tool's config:** `./dotme track ~/.config/<tool>` moves it into the
  repo, links it back and records it in `.dotme.lock` (`-m`/`-l` for
  macOS/Linux only).
- **Neovim plugins:** turn one off without touching its spec by adding its
  name to `vim.g.disabled_plugins` (read in `nvim/lua/config/lazy.lua`).
