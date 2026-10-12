function clone -d "bin/clone, then cd into the worktree"
    set -l dir (command clone $argv); and cd $dir
end
