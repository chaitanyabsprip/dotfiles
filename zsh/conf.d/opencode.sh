#!/bin/zsh

_have() { type "$1" >/dev/null 2>&1; }

_have opencode && prepend_path "$HOME/.opencode/bin"
