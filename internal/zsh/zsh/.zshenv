#!/bin/zsh

# Read by every zsh (login, interactive, scripts): keep it cheap.
export ZDOTDIR="$HOME/.config/zsh"
. "$HOME/dotfiles/shell/env.sh"
typeset -U path && path=($path) # dedupe, keeping the first occurrence
