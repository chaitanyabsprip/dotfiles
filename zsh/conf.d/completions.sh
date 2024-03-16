#!/bin/zsh

autoload -U compinit && compinit # Load and initialise completion system.
zstyle ':completion:*' menu select
_comp_options+=(globdots) # Include hidden files in completions.
