#!/bin/zsh

_have() { type "$1" >/dev/null 2>&1; }

_have jira && {
	# source <(jira completion zsh)
	alias issues='jira issue list -a$(jira me) -s~Done'
}
