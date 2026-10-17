#!/bin/zsh

# macOS /etc/zprofile runs path_helper after .zshenv, pushing /usr/bin ahead
# of Homebrew. Re-sourcing puts our PATH order back.
. "$DOTFILES/shell/env.sh"
path=($path) # -U from .zshenv only dedupes on array assignment
