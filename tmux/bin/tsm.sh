#!/bin/sh

_depends() { type "$1" >/dev/null 2>&1 ||
	{ echo "${0##*/} depends on $1, please install and try again." &&
		exit 1; }; }

_depends fzf-tmux

number_of_sessions() {
	tmux list-sessions |
		wc -l |
		sed "s/ //g"
}

current=$(tmux display-message -p "#S")

if [ "$(number_of_sessions)" -eq 1 ]; then
	tmux display-message "Only one session"
	exit 0
fi

if [ $# -eq 1 ]; then
	selected=$1
else
	selected=$(tmux ls -F "#{session_name}" | fzf-tmux -p 30%,30% --reverse)
fi

if [ "$current" = "$selected" ] || [ -z "$selected" ]; then
	exit 0
else
	tmux switch -t "$selected"
fi
