#!/bin/sh

number_of_sessions() {
	tmux list-sessions |
		wc -l |
		sed "s/ //g"
}

switch_session() {
	if [ "$(number_of_sessions)" -eq 2 ]; then
		tmux new-session -s home -c ~
	else
		tmux switch-client -l
	fi
}

current_session=$(tmux display -p "#{client_session}")
switch_session
tmux kill-session -t "$current_session"
