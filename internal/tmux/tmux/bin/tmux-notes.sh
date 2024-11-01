#!/bin/sh

_depends() { type "$1" >/dev/null 2>&1 ||
	{ echo "${0##*/} depends on $1, please install and try again." &&
		exit 1; }; }

_depends fzf-tmux

if [ $# -eq 1 ]; then
	selected=$1
else
	selected=$(
		fd .md "$NOTESPATH" |
			sed "s,$NOTESPATH/,," |
			fzf-tmux -p
	)
fi

if [ -z "$selected" ]; then
	exit 0
fi

selected="$(fd .md "$NOTESPATH" | grep -w "$selected")"

if tmux has-session -t notes 2>/dev/null; then
	tmux switch-client -t notes
	tmux new-window -c "$NOTESPATH" "zsh -lc 'nvim $selected'"
else
	tmux new-session -d -s notes -c "$NOTESPATH" "zsh -lc 'nvim $selected'; zsh"
	tmux switch-client -t notes
fi
