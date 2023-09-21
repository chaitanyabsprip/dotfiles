#!/bin/sh

_depends() { type "$1" >/dev/null 2>&1 ||
	{ echo "${0##*/} depends on $1, please install and try again." &&
		exit 1; }; }

_depends fzf-tmux

if [ $# -eq 1 ]; then
	selected=$1
else
	selected=$(
		(
			fd . -HIat d -d 1 -E .git \
				~/Programs \
				~/projects/apps \
				~/projects/forks \
				~/projects/languages \
				~/projects/languages/go \
				~/projects/languages/dart \
				~/projects/languages/lua \
				~/projects/languages/python
			fd . -HIt d -d 2 -E .git -E root \
				~/projects/apps/prep \
				~/projects/apps/work
			fd root -d 1 -t d \
				~/projects/apps/prep \
				~/projects/apps/work
			fd dotfiles ~ -at d -d 1
			fd notes ~/projects -at d -d 1
			fd . -HIad 1 -E .git ~/.config
		) | sort -u | fzf-tmux -p --border \
			--border-label=" Sessionizer " \
			--border-label-pos=6:bottom
	)
fi

if [ -z "$selected" ]; then
	exit 0
fi

selected_name=$(basename "$selected" | tr . _)
tmux_running=$(pgrep tmux)

if [ -z "$TMUX" ] && [ -z "$tmux_running" ]; then
	tmux new-session -s "$selected_name" -c "$selected"
	exit 0
fi

if ! tmux has-session -t="$selected_name" 2>/dev/null; then
	tmux new-session -ds "$selected_name" -c "$selected"
fi

tmux switch-client -t "$selected_name"
