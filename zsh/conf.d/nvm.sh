#!/bin/sh

lazynvm() {
	unset -f nvm node npm
	NVM_DIR="$HOME"/programs/nvm
	[ -s "$NVM_DIR" ] && export NVM_DIR="$NVM_DIR"
	[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh" # This loads nvm
}

nvm() {
	lazynvm
	nvm "$@"
}

node() {
	lazynvm
	node "$@"
}

npm() {
	lazynvm
	npm "$@"
}

nvim() {
	lazynvm
	nvim "$@"
}
