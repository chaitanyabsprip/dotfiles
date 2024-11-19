#!/bin/zsh

set -o vi
bindkey -v

bindkey -M viins '^b' vi-backward-blank-word
bindkey -M viins '^f' vi-forward-blank-word
bindkey -M viins '^a' beginning-of-line
bindkey -M viins '^e' end-of-line
