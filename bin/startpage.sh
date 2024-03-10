#!/bin/sh

_have() { type "$1" >/dev/null 2>&1; }

PORT=8910

setup() {
	echo "Hai bhai" >~/.temp
	cache_dir="${XDG_CACHE_HOME:-$HOME/.cache}"
	log_file="$cache_dir/minwall/server.log"
	cd ~/.config/minwall/ || exit 1
	[ ! -d "$cache_dir/minwall" ] && mkdir "$cache_dir/minwall"
	[ ! -f "$log_file" ] && touch "$log_file"
}

if _have browser-sync &&
	! { pgrep -f -l browser-sync | grep -q $PORT; }; then
	setup &&
		nohup browser-sync --host 127.0.0.1 --port $PORT >/dev/null 2>&1 &
fi

if _have http-server && ! { pgrep -fl http-server | grep -q $PORT; }; then
	setup &&
		nohup http-server -p $PORT >"$log_file" &
fi
