# Environment, aliases and fzf options are shared POSIX files written in a
# subset fish can source (see shell/env.sh). Only fish builtins live here.

source $HOME/dotfiles/shell/env.sh
set -gx PATH (printf '%s\n' $PATH | awk '!seen[$0]++') # dedupe, keeping the first occurrence

status is-interactive; or return

set -g fish_greeting
set -g fish_key_bindings fish_vi_key_bindings # extra binds: functions/fish_user_key_bindings.fish
source $HOME/.config/fish/themes/tokyonight_storm.fish

# Shared with bash and zsh
source $DOTFILES/shell/aliases
type -q fzf; and source $DOTFILES/shell/fzf.sh; and fzf --fish | source

type -q mise; and mise activate fish | source
type -q oh-my-posh; and oh-my-posh init fish --config $DOTFILES/oh-my-posh.rc.toml | source

# Drop-in config: fish/conf.d/*.fish is loaded by fish itself (local.fish is gitignored)
