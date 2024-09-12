#!/bin/sh

_have() { type "$1" >/dev/null 2>&1; }

_have pipx || return

append_path() {
	[ ! -d "$1" ] && return 1
	case ":$PATH:" in
	*:"$1":*) ;;
	::) export PATH="$1" ;;
	*) export PATH="$PATH:$1" ;;
	esac
}

export PIPX_BIN_DIR="$HOME"/.local/bin/.pipx
append_path "$PIPX_BIN_DIR"

eval "$(register-python-argcomplete pipx)"
