#!/bin/sh

{
	export PYENV_ROOT="$HOME/.pyenv"
	[ -d "$PYENV_ROOT"/bin ] && export PATH="$PYENV_ROOT/bin:$PATH" &&
		zsh-defer eval "$(pyenv init -)"
}
