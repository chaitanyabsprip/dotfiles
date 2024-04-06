#!/bin/sh

OLDBIND='c-a'
NEWBIND='c-b'
if [ -n "$SSH_TTY" ]; then
	tmux unbind-key "$OLDBIND"
	tmux set -g prefix "$NEWBIND"
	tmux bind "$NEWBIND" send-prefix
fi
