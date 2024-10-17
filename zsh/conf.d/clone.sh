#!/bin/sh

clone() {
	repo="$1"
	shift
	repo="$(echo "$repo" | sed 's,\(^https://github.com/|^git@github.com:\),,')"
	case "$repo" in
	*/*) user=$(echo "$repo" | cut -d'/' -f1) ;;
	*) user="${GITUSER:-Chaitanyabsprip}" ;;
	esac
	name=$(echo "$repo" | sed 's|.*/||')
	userd="${PROJECTS:-$HOME/projects}/$user"
	localPath="$userd/$name"
	[ -d "$localPath" ] && cd "$localPath" && return
	: $(mkdir -p "$userd")
	cd "$userd" || :
	echo gh repo clone "$user/$name" -- --recurse-submodule "$@"
	gh repo clone "$user/$name" -- --recurse-submodule "$@"
	cd "$name" || :
}
