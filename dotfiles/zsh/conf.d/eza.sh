#!/bin/zsh

_have() { type "$1" >/dev/null 2>&1; }

# Exit if the 'eza' command could not be found
_have eza && {

	# Use the --git flag if the installed version of eza supports git
	# Related to https://github.com/ogham/exa/issues/978
	if eza --version | grep -q '+git'; then
		alias ls='eza -F --group-directories-first --icons --git -1'
		alias ll='ls -lhF --git'
		alias la='ll -a'
	else
		alias ls='eza --group-directories-first --icons -1'
		alias ll='ls -lhF'
		alias la='ll -a'
	fi

	alias lt='eza --tree -ahD -L=2 --icons --git'
}
