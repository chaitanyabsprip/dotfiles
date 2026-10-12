#!/bin/sh

# Functions for bash and zsh. Only things that must change the current shell
# (cd, export, source) belong here; fish versions live in fish/functions/.

# md <dir>: mkdir + cd
md() { [ -n "$1" ] || return 1; mkdir -p "$1" && cd "$1"; }

# clone <user/repo|repo>: bin/clone does the work and prints where to cd
clone() { dir="$(command clone "$@")" && cd "$dir"; }

# jump [cmd]: fuzzy-pick a work dir, then `cmd` it (default: cd)
jump() {
	selected="$(SHORT=1 x work dirs | fzf --height=20% --border --margin=0 --padding=0)" || return
	${1:-cd} "$(x work dirs | grep -w "$selected$")"
}
