#!/bin/sh

_have() { type "$1" >/dev/null 2>&1; }

if _have browser-sync &&
	! { pgrep -f -l browser-sync | grep -q 8000; }; then
	cd ~/.config/minwall/ || exit 1
	nohup browser-sync --host 127.0.0.1 --port 8000 >/dev/null 2>&1 &
fi
