set -g message-command-style bg=default,fg=yellow
set -g message-style bg=default,fg=yellow

set -g status-interval 1
set -g status-position top
set -g status-style "fg=white,bg=default"
set -g automatic-rename on
set -g automatic-rename-format "\
#($HOME/.config/tmux/bin/tmux-icon-name.sh #{pane_current_command})"

set -g window-size latest
set -g status-left-length 100
set -g status-left '\
#[bg=#a3fcfe,fg=#1a1a1a,bold] #S \
#[bg=default]\
'

set -g status-right-length 100
set -g status-right '\
#(gitmux -cfg $HOME/.config/tmux/gitmux.conf "#{pane_current_path}")\
#[fg=green,bold]\
#{?#{!=:#(pomo), - }, #(pomo),}\
#[fg=white,bg=default,nobold]\
'

set -g window-status-format '#[fg=white]#{?window_index, ,}#{window_name} '
set -g window-status-current-format '#[bg=#161927,fg=#a3fcfe] #{window_name} '

set -g @minimal-statusline false
