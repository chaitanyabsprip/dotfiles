#!/bin/sh

clone() {
	repo="$1"
	shift
	repo="$(echo "$repo" | sed 's,\(^https://github.com/|^git@github.com:\),,')"
	case "$repo" in
	*/*) user="$(echo "$repo" | cut -d'/' -f1)" ;;
	*) user="${GITUSER:-Chaitanyabsprip}" ;;
	esac
	if [ "$user" = "${GITUSER:-Chaitanyabsprip}" ]; then
		user=""
	fi
	name=$(echo "$repo" | sed 's|.*/||')
	userd="${PROJECTS:-$HOME/projects}/$user"
	localPath="$(echo "$userd/$name" | sed 's|/$||')"
	[ -d "$localPath" ] && cd "$localPath" && return
	: "$(mkdir -p "$userd")"
	cd "$userd" || :
	echo gh repo clone "$user/$name" -- --recurse-submodule "$@"
	gh repo clone "$user/$name" -- --recurse-submodule "$@"
	cd "$name" || :
}
