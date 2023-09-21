#!/bin/sh
autoload -U compinit # Load and initialise completion system.
zstyle ':completion:*' menu select
compinit
_comp_options+=(globdots)	# Include hidden files in completions.
source <(/opt/homebrew/bin/jira completion zsh)
