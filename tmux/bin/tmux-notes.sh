#!/bin/sh

_depends() { type "$1" >/dev/null 2>&1 ||
	{ echo "${0##*/} depends on $1, please install and try again." &&
		exit 1; }; }

_depends fzf-tmux

if [ $# -eq 1 ]; then
	selected=$1
else
	selected=$(fd .md ~/projects/notes | fzf-tmux -p)
fi

if [ -z "$selected" ]; then
	exit 0
fi

if tmux has-session -t notes 2>/dev/null; then
	filename=$(basename "$selected")
	tmux switch-client -t notes
	tmux new-window -P -c ~/projects/notes -t "notes:$filename" "/bin/dash -c 'nvim $selected'; zsh"
else
	tmux new-session -d -s notes -c ~/projects/notes "/bin/dash -c 'nvim $selected'; zsh"
	tmux switch-client -t notes
fi
