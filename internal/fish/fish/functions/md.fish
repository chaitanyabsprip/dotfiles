function md -a dir -d "mkdir + cd"
    test -n "$dir"; or return 1
    mkdir -p $dir; and cd $dir
end
