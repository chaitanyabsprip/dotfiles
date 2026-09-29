#!/bin/zsh

[ -s "$HOME/.bum/bin/bum" ] && {
	export BUM_INSTALL="$HOME/.bum"
	export PATH="$BUM_INSTALL/bin:$PATH"
}

[ -s "$HOME/.bun/bin/bun" ] && {
	export BUN_INSTALL="$HOME/.bun"
	export PATH="$BUN_INSTALL/bin:$PATH"
}
