#!/bin/zsh

_activate_venv() {
	[ -d .venv ] && . "$(pdhas .venv)/bin/activate"
}

autoload -Uz add-zsh-hook
add-zsh-hook chpwd _activate_venv
