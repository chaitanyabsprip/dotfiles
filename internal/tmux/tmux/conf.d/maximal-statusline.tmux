set -g message-command-style "bg=#{@c_message_command_bg},fg=#{@c_message_command_fg}"
set -g message-style "bg=#{@c_message_bg},fg=#{@c_message_fg}"

set -g status-interval 1
set -g status-position top
set -g status-style "fg=#{@c_status_fg},bg=#{@c_status_bg}"
set -g automatic-rename on
set -g automatic-rename-format "\
#($HOME/.config/tmux/bin/tmux-icon-name.sh #{pane_current_command})"

set -g window-size latest
set -g status-left-length 100
set -g status-left '\
#[bg=#{@c_session_bg},fg=#{?client_prefix,#{@c_red},#{@c_session_bg}}]▎\
#[bg=#{@c_session_bg},fg=#{@c_session_fg},bold]#S \
#[bg=#{@c_status_bg},fg=#{@c_status_fg}]\
'

set -g status-right-length 100
set -g status-right '\
#(gitmux -cfg $HOME/.config/tmux/gitmux.conf "#{pane_current_path}")\
#[fg=#{@c_green},bold]\
#{?#{!=:#(pomo), - }, #(pomo),}\
#[fg=#{@c_status_fg},bg=#{@c_status_bg},nobold]\
'

set -g window-status-format '\
#[fg=#{@c_window_fg},bg=#{@c_window_bg}]#{?window_index, ,}#{window_name} \
'
set -g window-status-current-format '\
#[bg=#{@c_window_current_bg},fg=#{@c_window_current_fg}] #{window_name} \
'

set -g @minimal-statusline false
