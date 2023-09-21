cursor_mode() {
	# See https://ttssh2.osdn.jp/manual/4/en/usage/tips/vim.html for cursor shapes
	cursor_block='\e[2 q'
	cursor_beam='\e[6 q'

	function zle-keymap-select {
		if [[ ${KEYMAP} == vicmd ]] ||
			[[ $1 = 'block' ]]; then
			echo -ne $cursor_block
		elif [[ ${KEYMAP} == main ]] ||
			[[ ${KEYMAP} == viins ]] ||
			[[ ${KEYMAP} = '' ]] ||
			[[ $1 = 'beam' ]]; then
			echo -ne $cursor_beam
		fi
	}

	zle-line-init() { echo -ne $cursor_beam; }

	zle -N zle-keymap-select
	zle -N zle-line-init
}

cursor_mode

bindkey -M viins '^b' vi-backward-blank-word
bindkey -M viins '^f' vi-forward-blank-word
bindkey -M viins '^a' beginning-of-line
bindkey -M viins '^e' end-of-line
