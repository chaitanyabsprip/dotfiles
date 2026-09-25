set -g @resurrect-save 'M-C-s'
set -g @resurrect-restore 'M-C-r'
set -g @resurrect-hook-post-restore-all 'dot claude hook restore'
run '~/.config/tmux/bin/tmux-resurrect/resurrect.tmux'
