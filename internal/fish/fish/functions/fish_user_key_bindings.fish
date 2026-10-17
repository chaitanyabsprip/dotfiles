function fish_user_key_bindings
    bind -M insert \ca beginning-of-line
    bind -M insert \ce end-of-line
    bind -M insert \co 'commandline -r jump; commandline -f execute'
end
