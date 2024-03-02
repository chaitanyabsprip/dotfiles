#!/bin/zsh

fzf-history-widget-accept() {
	fzf-history-widget
	zle accept-line
}
zle -N fzf-history-widget-accept
bindkey '^X' fzf-history-widget-accept

# man page search widget for zsh
# fzf-man-widget() {
#   batman="man {1} | col -bx | bat --language=man --plain --color always --theme=\"Monokai Extended\""
#    man -k . | sort \
#    | awk -v cyan=$(tput setaf 6) -v blue=$(tput setaf 4) -v res=$(tput sgr0) -v bld=$(tput bold) '{ $1=cyan bld $1; $2=res blue;} 1' \
#    | fzf  \
#       -q "$1" \
#       --ansi \
#       --tiebreak=begin \
#       --prompt=' Man > '  \
#       --preview-window '50%,rounded,<50(up,85%,border-bottom)' \
#       --preview "${batman}" \
#       --bind "enter:execute(man {1})" \
#       --bind "alt-c:+change-preview(cht.sh {1})+change-prompt(ﯽ Cheat > )" \
#       --bind "alt-m:+change-preview(${batman})+change-prompt( Man > )" \
#       --bind "alt-t:+change-preview(tldr --color=always {1})+change-prompt(ﳁ TLDR > )"
#   zle reset-prompt
# }
# # `Ctrl-H` keybinding to launch the widget (this widget works only on zsh, don't
# # know how to do it on bash and fish (additionaly pressing`ctrl-backspace` will
# # trigger the widget to be executed too because both share the same keycode)
# bindkey '^h' fzf-man-widget
# zle -N fzf-man-widget
#
