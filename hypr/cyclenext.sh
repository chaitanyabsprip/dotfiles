#!/bin/sh

if hyprctl activewindow | grep -q grouped | grep -wq 'grouped: 0'; then
	if [ -z "$1" ]; then
		hyprctl dispatch cyclenext >/dev/null
	else
		hyprctl dispatch cyclenext prev >/dev/null
	fi
else
	if [ -z "$1" ]; then
		hyprctl dispatch changegroupactive f >/dev/null
	else
		hyprctl dispatch changegroupactive b >/dev/null
	fi
fi
