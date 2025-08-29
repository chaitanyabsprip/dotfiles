#!/bin/zsh

_have() { type "$1" >/dev/null 2>&1; }

_have uv && { eval "$(uv generate-shell-completion zsh)"; }
