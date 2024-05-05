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

		if is_tmux_running; then
			session_name="$(basename "$newpath")"
			tmux new-session -s "$session_name" -c "$newpath"
			exit 0
		fi

		if session_exists "$newpath"; then
			tmux switch-client -t "$(list_sessions | grep -w "$newpath" | cut -d '=' -f 1)"
			exit 0
		fi

		oldpath="$(list_sessions | grep -w "$session_name" | cut -d '=' -f 2)"
		read -r session_name old_session_new_name <<EOF
$(diffp "$newpath" "$oldpath")
EOF
		if session_exists "$session_name"; then
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
		workdirs | grep -w "$newpath"
	}

	is_tmux_running() {
		tmux_running=$(pgrep tmux)
		[ -z "$TMUX" ] && [ -z "$tmux_running" ]
	}

	session_exists() {
		list_sessions | grep -wq "$1"
	}

	list_sessions() {
		tmux ls -F '#{session_name}=#{session_path}'
	}

	diffp() {
		p1="$1"
		p2="$2"
		name1="$(basename "$p1")"
		name2="$(basename "$p2")"

		while [ "$name1" = "$name2" ] || [ -z "$name1" ] || [ -z "$name2" ]; do
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
