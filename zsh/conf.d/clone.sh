#!/bin/sh

# clone <user/repo|repo> [git-clone-args]
# Bare clone into ~/projects/<user>/<name>/.git, then one worktree per branch
# as sibling dirs. Home-base worktree is the remote default branch (main/master).
clone() {
	repo="$1"
	shift
	repo="$(echo "$repo" | sed -E 's,^(https://github.com/|git@github.com:),,; s,\.git$,,')"
	case "$repo" in
	*/*) user="$(echo "$repo" | cut -d'/' -f1)" ;;
	*) user="${GITUSER:-Chaitanyabsprip}" ;;
	esac
	name=$(echo "$repo" | sed 's|.*/||')
	userd="${PROJECTS:-$HOME/projects}/$user"
	if [ "$user" = "${GITUSER:-Chaitanyabsprip}" ]; then
		userd="${PROJECTS:-$HOME/projects}"
	fi
	localPath="$userd/$name"
	[ -d "$localPath" ] && cd "$localPath" && return
	mkdir -p "$localPath" && cd "$localPath" || return

	echo gh repo clone "$user/$name" .git -- --bare "$@"
	gh repo clone "$user/$name" .git -- --bare "$@" || return
	# bare clones omit the fetch refspec, so worktree add can't auto-track origin/*
	git config remote.origin.fetch "+refs/heads/*:refs/remotes/origin/*"
	git fetch origin
	default="$(git symbolic-ref --short HEAD)"
	# drop the untracked local copies of every remote branch; switch/worktree add
	# then recreates them on demand with upstream set
	git for-each-ref refs/heads --format='%(refname:short)' | grep -vx "$default" | xargs -r git branch -D
	git branch --set-upstream-to="origin/$default" "$default"
	git worktree add "$default"
	cd "$default" || :
}
