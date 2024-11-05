#!/bin/sh

onevent() {
	echo "$*"
	case "$*" in
	*) ;;
	esac
}

socat -U - UNIX-CONNECT:/tmp/hypr/"$HYPRLAND_INSTANCE_SIGNATURE"/.socket2.sock | while read -r line; do onevent "$line"; done
