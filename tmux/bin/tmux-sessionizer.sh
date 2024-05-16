#!/bin/sh

sessionizer() {
	_depends() { type "$1" >/dev/null 2>&1 ||
		{ echo "${0##*/} depends on $1, please install and try again." &&
			exit 1; }; }

	_depends fzf-tmux

	main() {
		newpath="$(select_path "$@")"
		[ -z "$newpath" ] && exit 0
		session_name="$(basename "$newpath")"

		if tmux_inactive; then
			session_name="$(basename "$newpath")"
			exec tmux new-session -s "$session_name" -c "$newpath"
		fi

		if sessions_has_match "=$newpath$"; then
			exec tmux switch-client -t "$(list_sessions | grep "=$newpath$" | cut -d '=' -f 1)"
		fi

		oldpath="$(list_sessions | grep -w "$session_name" | cut -d '=' -f 2)"
		read -r session_name old_session_new_name <<EOF
$(diffp "$newpath" "$oldpath")
EOF
		if sessions_has_match "$session_name"; then
			tmux rename-session -t "$session_name" "$old_session_new_name"
		fi

		tmux new-session -ds "$session_name" -c "$newpath"
		tmux switch-client -t "$session_name"
	}

	select_path() {
		if [ -n "$1" ]; then
			newpath=$1
		else
			newpath=$(
				workdirs -s | fzf-tmux -p 30% --border \
					--border-label=' Sessionizer ' \
					--border-label-pos=6:bottom
			)
		fi
		[ -z "$newpath" ] && exit 0
		workdirs | grep -w "$newpath$"
	}

	tmux_inactive() {
		tmux_running=$(pgrep tmux)
		[ -z "$TMUX" ] && [ -z "$tmux_running" ]
	}

	sessions_has_match() {
		list_sessions | grep -q "$1"
	}

	list_sessions() {
		tmux ls -F '#{session_name}=#{session_path}' | grep -v '^scratch='
	}

	diffp() {
		p1="$1"
		p2="$2"
		name1="$(basename "$p1")"
		name2="$(basename "$p2")"
		while [ "$name1" = "$name2" ] && [ -n "$name1" ] && [ -n "$name2" ]; do
			p1=$(dirname "$p1")
			p2=$(dirname "$p2")
			name1="$(basename "$p1")/$(basename "$name1")"
			name2="$(basename "$p2")/$(basename "$name2")"
		done
		echo "$name1 $name2"
	}

	main "$@"
}

sessionizer "$@"
