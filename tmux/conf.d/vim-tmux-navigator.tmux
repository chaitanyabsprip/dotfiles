# Smart pane switching with awareness of Vim splits.
# See: https://github.com/christoomey/vim-tmux-navigator
is_vim="ps -o state= -o comm= -t '#{pane_tty}' \
    | grep -iqE '^[^TXZ ]+ +(\\S+\\/)?g?(view|l?n?vim?x?|fzf)(diff)?$'"
bind-key -n -N "Jump to pane on left"   C-h if-shell "$is_vim" 'send-keys C-h'  'select-pane -L'
bind-key -n -N "Jump to pane below"     C-j if-shell "$is_vim" 'send-keys C-j'  'select-pane -D'
bind-key -n -N "Jump to pane above"     C-k if-shell "$is_vim" 'send-keys C-k'  'select-pane -U'
bind-key -n -N "Jump to pane on right"  C-l if-shell "$is_vim" 'send-keys C-l'  'select-pane -R'
tmux_version='$(tmux -V | sed -En "s/^tmux ([0-9]+(.[0-9]+)?).*/\1/p")'
if-shell -b '[ "$(echo "$tmux_version < 3.0" | bc)" = 1 ]' \
    "bind-key -n 'C-\\' if-shell \"$is_vim\" 'send-keys C-\\'  'select-pane -l'"
if-shell -b '[ "$(echo "$tmux_version >= 3.0" | bc)" = 1 ]' \
    "bind-key -n 'C-\\' if-shell \"$is_vim\" 'send-keys C-\\\\'  'select-pane -l'"

bind-key -T copy-mode-vi -N "Copy mode vi: Jump to pane on left"   C-h select-pane -L
bind-key -T copy-mode-vi -N "Copy mode vi: Jump to pane below"     C-j select-pane -D
bind-key -T copy-mode-vi -N "Copy mode vi: Jump to pane above"     C-k select-pane -U
bind-key -T copy-mode-vi -N "Copy mode vi: Jump to pane on right"  C-l select-pane -R
bind-key -T copy-mode-vi -N "Copy mode vi: Jump to last pane"      C-\ select-pane -l
