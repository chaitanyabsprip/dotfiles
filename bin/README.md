# bin

Personal scripts, on `PATH`. Tags: **[linux]** Linux-only, **[hypr]**
Hyprland/Wayland desktop, **[gnu]** needs GNU coreutils/date.

## git

- `copy-branch` — branch off HEAD as `<current>-<suffix>` and switch to it
- `g-` — `git switch -`
- `gbd` — fzf-tmux pick branches to delete (`-f` force)
- `gco` — fzf-tmux pick a branch or tag to check out
- `gd` — `git diff`
- `gf` — `git fetch`
- `gissue` — compact `gh issue list` with labels
- `git-when` — per file in cwd: last commit, or untracked/staged/ignored
- `gitcolor` — awk filter colouring `git status --porcelain` output
- `glog` — fzf git log browser with commit preview
- `gp` / `gpl` — `git push` / `git pull`
- `gs` — rich `git status` with dividers
- `gst` — `git status -sb`
- `gwa` — create a worktree named from a `kind/module/ticket/desc` branch
- `origins` — print origin remote URLs
- `pdgit` — print nearest ancestor dir containing `.git`
- `status` — porcelain status for every git repo under a dir
- `clone` — bare-clone `<user/repo>` into `~/projects` with a worktree per branch; prints the dir (the `clone` shell function cds into it)

## tmux

- `cheat` — fzf pick a language/command and query cht.sh
- `claude-tmux` — record Claude Code sessions per pane so tmux-resurrect can resume them
- `gitpeek` — git log in a tmux popup
- `gitwatch` — live-updating git log in a split pane
- `tmux-peek` — show the tmux status bar for 5 seconds
- `tssh` — run tmux on a separate socket inside an ssh session

## install

- `distro` — print `<os-release ID>:<VERSION_ID>` (e.g. `arch:`, `debian:12`, `Darwin:25.5.0`)
- `install-android-sdk` — Android SDK, Java 17 and Gradle **[linux]**
- `install-bat` / `install-eza` / `install-fd` / `install-gitmux` / `install-gitui` / `install-jq` / `install-rg` / `install-yq` — install via package manager, else latest GitHub release into `$SCRIPTS`
- `install-flutter` — Flutter via fvm plus build deps **[linux]**
- `update-neovim` — build and install neovim from source (default `nightly`)

## go

- `goclean` — clear Go build and module caches
- `godistbuild` — build a package for every `GOOS/GOARCH`
- `gor` — `go run` main.go/generate.go/args
- `got` — re-run `go test` on file changes

## system

- `b` — brightness via brightnessctl (`i`/`d`/`max`/`m`in/`f`ull/0-10) **[linux]**
- `battery` — battery state from upower **[linux]**
- `cdpath` / `path` — print `$CDPATH` / `$PATH` one per line
- `cleanlinks` — delete broken symlinks in cwd
- `dev` — docker container with a dir mounted at `/root/<name>`
- `ffrestore` / `ffsession` — save and reopen Firefox session URLs
- `islinux` / `ismac` — exit 0 on that OS
- `listening` — list listening ports
- `myip` — external IP via DNS
- `nopath` — run a command with empty `PATH`
- `notify` — desktop notification (terminal-notifier, say, or notify-send)
- `opacity` / `fontsize` — edit Alacritty opacity / font size
- `powermenu` — wofi power menu (shutdown, reboot, logout, suspend) **[linux]** **[hypr]**
- `pubkey` — print an ssh public key (default id_rsa/id_ed25519)
- `run` — throwaway privileged docker container for an image
- `rot` — run a command whenever a trigger file is touched (inotify) **[linux]**
- `scann` — nmap all open ports then service-scan them
- `sshk` — generate/list/print ed25519 ssh keys
- `startpage.sh` — serve the minwall start page with python http.server **[linux]**
- `sw` — set wallpaper via hyprpaper or awww **[linux]** **[hypr]**
- `tcpu` / `tmem` — top 10 processes by CPU / memory
- `topdu` — largest entries in cwd
- `v` — volume via wpctl **[linux]**
- `withenv` — run a command with a `.env` file exported

## text and files

- `align` / `ralign` — align stdin on a character / right-align
- `bold` / `italic` — print text bold / italic
- `chmox` — `chmod +x`
- `edit` — fzf pick a file (6 levels deep) and open it in nvim
- `compress` / `extract` — create / unpack archives by extension
- `countf` — count files under a path, optionally by name glob
- `diffc` — diff the output of two commands
- `fmttable` — format markdown tables from stdin
- `istext` — print a file's MIME top-level type
- `lslinks` — list symlinks under cwd
- `md` — markdown helpers (`toc`, `imagelinks`)
- `minify` — collapse stdin to one line
- `newest` — newest file in a dir **[gnu]**
- `outline1` — keep only lines starting with `*`
- `shorten-pwd` — `$PWD` with abbreviated parent dirs
- `symlink` — `ln -sf` with arguments in `<from> <to>` order
- `urlencode` / `urldecode` — URL-encode / decode

## shell and scripting helpers

- `changed` / `haschanged` / `onchange` — detect file changes since last run, run a command on change
- `confirm` — `[y/N]` prompt, exit status is the answer
- `depends` — exit (and kill the caller) if a command is missing
- `error` / `heading` / `inprogress` / `success` / `warning` — styled status lines
- `have` — exit 0 if all commands are on `PATH`
- `hem` — save a command line as a new script in `$SCRIPTS`
- `newx` — create or edit a script in `$SCRIPTS`
- `pdhas` — search upward for the nearest ancestor containing a file/dir
- `runme` — run a piped markdown sh codeblock and append its output
- `sc` — show a ShellCheck wiki page for an error code
- `shellprof` — profile which output lines of a command took longest
- `spinner` — show a spinner while a command runs
- `vic` — open a command on `PATH` in `$EDITOR`

## clipboard

- `pbc` / `pbp` — copy / paste (pbcopy, xclip or wl-copy)
- `pbr` — run clipboard contents as a script
- `csfp` — write clipboard to a file and make it executable

## notes, time and misc

- `back` / `br` — pomodoro work / break timers via pomo
- `daily` — open today's note, carrying over open todos
- `dic` — dict.org definition
- `emails` — list or fzf-pick saved email addresses
- `fltr` — `flutter run` with hot reload on save (entr)
- `lastdown` / `mvlast` — newest download / move it to cwd
- `logma` — live-view a file, refreshing on change (entr)
- `lorep` / `lores` — lorem paragraphs / sentences
- `meet` / `sheet` — new Google Meet / Sheet in the browser
- `now` / `till` — formatted date / seconds until a date **[gnu]**
- `ns` — dump the last N headings from the notes file
- `prj` — list `~/projects` **[gnu]**
- `remind` — notify after N seconds
- `weather` — wttr.in forecast
