#!/bin/sh

_have() { type "$1" >/dev/null 2>&1; }

PORT=8910

cache_dir="${XDG_CACHE_HOME:-$HOME/.cache}"
log_file="$cache_dir/minwall/server.log"

setup() {
	cd ~/.config/minwall/ || exit 1
	[ ! -d "$cache_dir/minwall" ] && mkdir "$cache_dir/minwall"
	[ ! -f "$log_file" ] && touch "$log_file"
}

started="false"
if _have python3 && ! { pgrep -fa http.server | grep -q $PORT; }; then
	started="python3 server"
	setup
	nohup /usr/bin/python3 -m http.server $PORT >>"$log_file" 2>&1 &
fi

if [ "$started" = "false" ]; then
	echo "Found neither browser-sync nor http-server or already running" >>"$log_file" 2>&1
else
	echo "Started $started ($(pgrep -fa http.server))" >>"$log_file" 2>&1
fi
