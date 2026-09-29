#!/bin/zsh

_have() { type "$1" >/dev/null 2>&1; }
_ismac() { [ "$(uname)" = "Darwin" ]; }

_ismac && prepend_path "$HOME/.opencode/bin"
