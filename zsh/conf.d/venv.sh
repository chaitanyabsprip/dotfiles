#!/bin/zsh

_activate_venv() {
	if [ -d ".venv" ]; then
		typeset -f python >/dev/null && lazypyenv
		source "$(pdhas .venv)/bin/activate"
	fi
}

add-zsh-hook chpwd _activate_venv
