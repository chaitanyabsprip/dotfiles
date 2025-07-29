#!/bin/zsh

HISTORY_IGNORE="(la|ls|cd|pwd|exit|history|cd -|cd ..|q|c|vim|nvim)"
SAVEHIST=200000
HISTFILE=~/.config/zsh/.zsh_history
HISTFILESIZE=1000000000
HISTSIZE=1000000000
ZLE_RPROMPT_INDENT=0

_have() { type "$1" >/dev/null 2>&1; }

# Load separated config files
plugin_file=$HOME/.config/zsh/plugs.zsh
config_dir=$HOME/.config/zsh/conf.d

[ -r $plugin_file ] && source $plugin_file

for conf in "$config_dir/"*.sh; do
	[ -r "$conf" ] && source "$conf"
done
unset conf

# key-bindings
autoload -Uz edit-command-line
zle -N edit-command-line
bindkey '^v' edit-command-line
bindkey -s ^o '^ujump\n'

# initialisations
_have rbenv && eval "$(rbenv init - zsh)"
# have note && zsh-defer eval "$(note completion zsh)"

# bun completions
[ -s "$HOME/.bun/_bun" ] && source "$HOME/.bun/_bun"
: # shell should start with a zero status code
