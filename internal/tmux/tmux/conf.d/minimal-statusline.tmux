set -g message-command-style "bg=#{@c_message_command_bg},fg=#{@c_message_command_fg}"
set -g message-style "bg=#{@c_message_bg},fg=#{@c_message_fg}"

set -g status-interval 1
set -g status-position top
set -g status-style "fg=#{@c_status_fg},bg=#{@c_status_bg}"
set -g automatic-rename on
set -g automatic-rename-format "\
#(dot tmux x icon #{pane_current_command})"

set -g window-size latest
set -g status-left-length 100
set -g status-left '\
#[fg=#{@c_status_fg},bg=#{@c_status_bg},nobold]\
'

set -g status-right-length 100
set -g status-right '\
#[fg=#{@c_session_bg},bg=#{@c_session_fg},bold] #S \
'

set -g window-status-format '#[fg=#{@c_window_fg}]#{?window_index, ,}#{window_name} '
set -g window-status-current-format '#[bg=#{@c_window_current_bg},fg=#{@c_window_current_fg}] #{window_name} '

set -g @minimal-statusline true
