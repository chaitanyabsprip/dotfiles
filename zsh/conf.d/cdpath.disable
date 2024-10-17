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
cdpath_add $HOME/projects
cdpath_add $HOME/.config
cdpath_add $HOME/programs
cdpath_add $HOME/.local/share/nvim/lazy
unfunction cdpath_add
