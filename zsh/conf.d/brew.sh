#!/bin/sh

have brew && {
	lazybrew() {
		unset -f brew
		eval "$(brew shellenv)"
	}

	brew() {
		lazybrew
		brew "$@"
	}
}
