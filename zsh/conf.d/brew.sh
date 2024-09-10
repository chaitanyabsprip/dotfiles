#!/bin/zsh

_have() { type "$1" >/dev/null 2>&1; }

_have brew && {
	lazybrew() {
		unset -f brew
		eval "$(brew shellenv)"
	}

	brew() {
		lazybrew
		brew "$@"
	}
}
