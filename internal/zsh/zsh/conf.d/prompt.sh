#!/bin/sh

if _have oh-my-posh; then
	eval "$(oh-my-posh init zsh -c ~/.config/oh-my-posh/oh-my-posh.rc.toml)"
elif _have starship; then
	eval "$(starship init zsh)"
fi
