#!/bin/zsh

export HISTORY_IGNORE="(la|ls|cd|pwd|exit|history|cd -|cd ..|q|c|vim|nvim)"
export HISTSIZE=200000

# Cannot use autocd option along with CDPATH
# setopt autocd                 # Automatically cd into typed directory.

setopt hist_expire_dups_first # delete duplicates first when HISTFILE size exceeds HISTSIZE
setopt hist_ignore_dups       # ignore duplicated commands history list
setopt hist_ignore_all_dups   # Delete old recorded entry if new entry is a duplicate
setopt hist_ignore_space      # ignore commands that start with space
setopt hist_verify            # show command with history expansion to user before running it
setopt inc_append_history     # add commands to HISTFILE in order of execution
stty stop undef               # Disable ctrl-s to freeze terminal.

set -o vi
bindkey -v

# _have() { type "$1" &>/dev/null; }

# Load seperated config files
plugin_file=$HOME/.config/zsh/plugs.zsh
config_dir=$HOME/.config/zsh/conf.d

[ -r $plugin_file ] && source $plugin_file

for conf in "$config_dir/"*.sh; do
	[ -r "$conf" ] && zsh-defer source "$conf"
done
unset conf

# key-bindings
have fd alias fzfdir="fd . -t d --max-depth 1 |\
  fzf  --height=40% --border --margin=0 --padding=0"
have fd && have fzf && bindkey -s '^f' '^ucd "$(fzfdir)"\n' # cd to folder with fzf
zsh-defer autoload -Uz edit-command-line
zle -N edit-command-line
zsh-defer bindkey '^v' edit-command-line
have jira && zsh-defer bindkey -s '^k' '^ujim\n'

# initialisations
have brew && zsh-defer eval "$(brew shellenv | grep -v 'export PATH')"
{
	export PYENV_ROOT="$HOME/.pyenv"
	[[ -d $PYENV_ROOT/bin ]] && export PATH="$PYENV_ROOT/bin:$PATH"
	zsh-defer eval "$(pyenv init -)"
}
have rbenv && zsh-defer eval "$(rbenv init - zsh)"
{
	NVM_DIR="$HOME"/programs/nvm
	[ -s $NVM_DIR ] && export NVM_DIR=$NVM_DIR
	[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm
}
have starship && eval "$(starship init zsh)"

# zsh-defer unfunction _have

[[ $TMUX ]] || tmux new -As home
