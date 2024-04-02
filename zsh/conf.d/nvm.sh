#!/bin/sh

{
	NVM_DIR="$HOME"/programs/nvm
	[ -s "$NVM_DIR" ] && export NVM_DIR="$NVM_DIR"
	[ -s "$NVM_DIR/nvm.sh" ] && zsh-defer \. "$NVM_DIR/nvm.sh" # This loads nvm
}
