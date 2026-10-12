#!/bin/zsh

# Interactive zsh. Environment comes from .zshenv (shell/env.sh).

HISTFILE="$ZDOTDIR/.zsh_history"
HISTSIZE=1000000
SAVEHIST=1000000
HISTORY_IGNORE="(la|ls|cd|pwd|exit|history|cd -|cd ..|q|c|vim|nvim)"
ZLE_RPROMPT_INDENT=0

_have() { (( $+commands[$1] )); }

# Shared with bash and fish
. "$DOTFILES/shell/aliases"
. "$DOTFILES/shell/functions.sh"
_have fzf && . "$DOTFILES/shell/fzf.sh"

# zsh only. Disable a conf.d file by renaming it to *.off; work.zsh is gitignored.
[ -r "$ZDOTDIR/plugs.zsh" ] && . "$ZDOTDIR/plugs.zsh"
for conf in "$ZDOTDIR"/conf.d/*.zsh(N); do . "$conf"; done
unset conf

: # shell should start with a zero status code
