#!/bin/sh

# Created by Zap installer
[ -f "${XDG_DATA_HOME:-$HOME/.local/share}/zap/zap.zsh" ] && source "${XDG_DATA_HOME:-$HOME/.local/share}/zap/zap.zsh"
plug "romkatv/zsh-defer"
zsh-defer plug "zsh-users/zsh-autosuggestions"
zsh-defer bindkey '^ ' autosuggest-accept
zsh-defer plug "zdharma-continuum/fast-syntax-highlighting" # Load syntax highlighting; should be last.

zsh-defer plug "zsh-users/zsh-history-substring-search"
zsh-defer bindkey -M vicmd 'k' history-substring-search-up
zsh-defer bindkey -M vicmd 'j' history-substring-search-down
zsh-defer bindkey "^P" history-substring-search-up
zsh-defer bindkey "^N" history-substring-search-down
zsh-defer bindkey "^[[A" history-substring-search-up
zsh-defer bindkey "^[[B" history-substring-search-down
