#!/bin/zsh

_have() { type "$1" >/dev/null 2>&1; }

_have pyenv && {
	lazypyenv() {
		unset -f python pip python3 pip3 nvim
		export PYENV_ROOT="$HOME/.pyenv"
		[ -d "$PYENV_ROOT"/bin ] && export PATH="$PYENV_ROOT/bin:$PATH" &&
			eval "$(pyenv init -)"
	}

	python() {
		lazypyenv
		python "$@"
	}

	pip() {
		lazypyenv
		pip "$@"
	}

	python3() {
		lazypyenv
		python3 "$@"
	}

	pip3() {
		lazypyenv
		pip3 "$@"
	}

	nvim() {
		lazypyenv
		nvim "$@"
	}
}
