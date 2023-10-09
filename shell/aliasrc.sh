#!/bin/sh

_have() { type "$1" >/dev/null 2>&1; }

# Only to be used as an alias to disable commands
donothing() { echo "$*" >/dev/null 2>&1; }

_have nvim && alias nvim='donothing'
# alias clear='donothing'
# alias exit='donothing'

_have pbcopy && alias pbc='pbcopy'
_have pbpaste && alias pbp='pbpaste'

alias mkdir='mkdir -pv'
alias chmox='chmod +x'
alias view='vi -R'

_have note && alias bm='note -b'
_have note && alias did='note -c'
_have note && alias todo='note -t'

_have gitui && alias gitui='gitui -t mocha.ron'
_have /usr/bin/vim && alias vi=/usr/bin/vim
_have nvim && alias vim=/usr/local/bin/nvim vimdiff='nvim -d'
_have bat && alias cat='bat'
_have bard-cli && alias \?='bard-cli -c ~/.config/zsh/.bardcli.yaml'
_have jira && alias jim="jira issue list \
  -a\$(jira me) \
  -q\"statuscategory in (New,'In Progress')\" \
  --plain | \
  tail -n +2 | \
  fzf-tmux \
  --preview=\"echo {} |\
      grep -Eo 'YOC-[0-9]+' |\
      xargs jira issue view |\
      tail -n +4\" \
  --preview-window=right,60%"

md() { [ -z "$1" ] && exit 1 || mkdir "$1" && cd "$1" || return; }

_have fd && _have fzf && edit() {
	fd . -Ht f -d 6 -E .git -E .DS_Store |
		fzf --height=20% --border --margin=0 --padding=0 \
			--bind='enter:become(nvim {})'
}
