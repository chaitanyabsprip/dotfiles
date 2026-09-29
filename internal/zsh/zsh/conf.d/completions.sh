#!/bin/zsh

fpath=(${ZDOTDIR:-$HOME/.config/zsh}/functions $fpath)

autoload -U compinit && compinit         # Load and initialise completion system.
autoload -U bashcompinit && bashcompinit # Load bash completions.
zstyle ':completion:*' menu select
zstyle ':completion:*:*:make:*' tag-order 'targets variables'
_comp_options+=(globdots) # Include hidden files in completions.

complete -C dot dot
complete -C x x
complete -C kimono kimono
complete -C tug tug
