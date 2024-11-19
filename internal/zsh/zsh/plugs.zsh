#!/bin/sh

# Created by Zap installer
[ ! -f "${XDG_DATA_HOME:-$HOME/.local/share}/zap/zap.zsh" ] && return
source "${XDG_DATA_HOME:-$HOME/.local/share}/zap/zap.zsh"
plug "zsh-users/zsh-autosuggestions"
bindkey '^ ' autosuggest-accept
plug "zdharma-continuum/fast-syntax-highlighting" # Load syntax highlighting; should be last.

plug "zsh-users/zsh-history-substring-search"
bindkey -M vicmd 'k' history-substring-search-up
bindkey -M vicmd 'j' history-substring-search-down
bindkey "^P" history-substring-search-up
bindkey "^N" history-substring-search-down
bindkey "^[[A" history-substring-search-up
bindkey "^[[B" history-substring-search-down
