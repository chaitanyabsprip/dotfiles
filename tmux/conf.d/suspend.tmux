set -g @suspend_key 'M-s'
set -g @suspend_suspended_options " \
status-position::bottom, \
status-style::bg=default, \
status-right:: ##S , \
status-right-style::fg=#{@c_green}, \
window-status-current-format::, \
window-status-format::, \
status-left::, \
"
run-shell ~/.config/tmux/bin/tmux-suspend/suspend.tmux
