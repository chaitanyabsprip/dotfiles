#!/bin/zsh

_have() { type "$1" >/dev/null 2>&1; }

_have fzf && {
	# fzf-history-widget-accept() {
	# 	fzf-history-widget
	# 	zle accept-line
	# }
	# zle -N fzf-history-widget-accept
	# bindkey '^X' fzf-history-widget-accept

	source <(fzf --zsh)
}
