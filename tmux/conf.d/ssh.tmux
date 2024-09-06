%if '#{||:#{==:#{socket_path},/tmp/tmux-1000/ssh},#{==:#{socket_path},/private/tmp/tmux-501/ssh}}'
unbind C-a
set -g prefix C-b
bind C-b send-prefix
%endif
run ''
