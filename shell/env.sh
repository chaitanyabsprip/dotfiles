#!/bin/sh

# Environment shared by sh, bash, zsh and fish.
#
# COMMON SUBSET: fish sources this file too, so stick to
#   export NAME="value"          (no bare NAME=value)
#   [ test ] && export ...       (no if/case/for, no functions)
#   "$(cmd)", "$VAR"             (no ${VAR:-default})
#   . file
#
# Safe to source repeatedly: every shell dedupes PATH afterwards, keeping the
# first occurrence, so re-prepending an entry moves it to the front.

export DOTFILES="$HOME/dotfiles"
export PROJECTS="$HOME/projects"
export NOTESPATH="$HOME/projects/library"
export DOWNLOADS="$HOME/downloads"
export SCRIPTS="$HOME/.local/bin"
export GITUSER="Chaitanyabsprip"

export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$HOME/.local/share"
export XDG_STATE_HOME="$HOME/.local/state"
export XDG_CACHE_HOME="$HOME/.local/cache"

export LANG="en_US.UTF-8"
export LC_ALL="en_US.UTF-8"

export VISUAL="nvim"
export EDITOR="nvim"
export BAT_THEME="material-rose-pine-dark"

# Toolchains
export GOPATH="$HOME/programs/go"
export GOBIN="$GOPATH/bin"
export FVM_CACHE_PATH="$HOME/programs/flutter"
export FLUTTER_SDK_ROOT="$FVM_CACHE_PATH/default"
export ANDROID_SDK_ROOT="$HOME/programs/android"
[ -d "$HOME/Library/Android/Sdk" ] && export ANDROID_SDK_ROOT="$HOME/Library/Android/Sdk"
export ANDROID_HOME="$ANDROID_SDK_ROOT"
[ -d /opt/homebrew/opt/llvm ] && export LDFLAGS="-L/opt/homebrew/opt/llvm/lib"
[ -d /opt/homebrew/opt/llvm ] && export CPPFLAGS="-I/opt/homebrew/opt/llvm/include"

# Subprocesses: computed once, then inherited by child shells.
[ -z "$JAVA_HOME" ] && [ -d /usr/lib/jvm/java-17-openjdk-amd64 ] && export JAVA_HOME="/usr/lib/jvm/java-17-openjdk-amd64"
[ -z "$JAVA_HOME" ] && [ -d /usr/lib/jvm/java-17-openjdk ] && export JAVA_HOME="/usr/lib/jvm/java-17-openjdk"
[ -z "$JAVA_HOME" ] && [ -x /usr/libexec/java_home ] && export JAVA_HOME="$(/usr/libexec/java_home 2>/dev/null)"
[ -z "$JIRA_API_TOKEN" ] && [ -r "$HOME/.secrets/jira-token" ] && export JIRA_API_TOKEN="$(cat "$HOME/.secrets/jira-token")"

# PATH: appended (fallbacks, last resort)
[ -d "$FLUTTER_SDK_ROOT/bin" ] && export PATH="$PATH:$FLUTTER_SDK_ROOT/bin"
[ -d "$FLUTTER_SDK_ROOT/bin/cache/dart-sdk/bin" ] && export PATH="$PATH:$FLUTTER_SDK_ROOT/bin/cache/dart-sdk/bin"
[ -d "$ANDROID_HOME/cmdline-tools/bin" ] && export PATH="$PATH:$ANDROID_HOME/cmdline-tools/bin"
[ -d "$ANDROID_SDK_ROOT/emulator" ] && export PATH="$PATH:$ANDROID_SDK_ROOT/emulator"
[ -d "$ANDROID_SDK_ROOT/platform-tools" ] && export PATH="$PATH:$ANDROID_SDK_ROOT/platform-tools"
[ -d "$ANDROID_SDK_ROOT/tools" ] && export PATH="$PATH:$ANDROID_SDK_ROOT/tools"
[ -d "$ANDROID_SDK_ROOT/tools/bin" ] && export PATH="$PATH:$ANDROID_SDK_ROOT/tools/bin"
[ -d "$HOME/.poetry/bin" ] && export PATH="$PATH:$HOME/.poetry/bin"
[ -d "$HOME/.cargo/bin" ] && export PATH="$PATH:$HOME/.cargo/bin"
[ -d "$HOME/.pub-cache/bin" ] && export PATH="$PATH:$HOME/.pub-cache/bin"
[ -d "$GOBIN" ] && export PATH="$PATH:$GOBIN"
[ -d "$XDG_DATA_HOME/nvim/mason/bin" ] && export PATH="$PATH:$XDG_DATA_HOME/nvim/mason/bin"
[ -d "$HOME/.antigravity/antigravity/bin" ] && export PATH="$PATH:$HOME/.antigravity/antigravity/bin"
[ -d "$HOME/.maestro/bin" ] && export PATH="$PATH:$HOME/.maestro/bin"

# PATH: prepended (beat the system), so the LAST line here ends up first
[ -d /opt/homebrew/sbin ] && export PATH="/opt/homebrew/sbin:$PATH"
[ -d /opt/homebrew/bin ] && export PATH="/opt/homebrew/bin:$PATH"
[ -d /opt/homebrew/opt/coreutils/libexec/gnubin ] && export PATH="/opt/homebrew/opt/coreutils/libexec/gnubin:$PATH"
[ -d "$HOME/.bun/bin" ] && export BUN_INSTALL="$HOME/.bun"
[ -d "$HOME/.bun/bin" ] && export PATH="$HOME/.bun/bin:$PATH"
[ -d "$HOME/.bum/bin" ] && export BUM_INSTALL="$HOME/.bum"
[ -d "$HOME/.bum/bin" ] && export PATH="$HOME/.bum/bin:$PATH"
[ -d "$HOME/.opencode/bin" ] && export PATH="$HOME/.opencode/bin:$PATH"
[ -d "$SCRIPTS" ] && export PATH="$SCRIPTS:$PATH"
[ -d "$SCRIPTS/overrides" ] && export PATH="$SCRIPTS/overrides:$PATH"
