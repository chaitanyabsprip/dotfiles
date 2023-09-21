#!/bin/zsh

cdpath_add() {
	case ":$CDPATH:" in
	*:$1:*) ;;
	*)
		[ -d "$1" ] && export CDPATH="$CDPATH:$1"
		;;
	esac
}

CDPATH="."
cdpath_add $HOME/dotfiles
cdpath_add $HOME/.config
cdpath_add $HOME/Programs
cdpath_add $HOME/projects
cdpath_add $HOME/projects/apps
cdpath_add $HOME/projects/apps/work
cdpath_add $HOME/projects/languages
cdpath_add $HOME/projects/forks
cdpath_add $HOME/.local/share/nvim/lazy
unfunction cdpath_add
