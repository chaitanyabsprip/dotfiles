function jump -a program -d "fuzzy-pick a work dir, then `program` it (default: cd)"
    set -q program[1]; or set program cd
    set -l selected (SHORT=1 x work dirs | fzf --height=20% --border --margin=0 --padding=0); or return
    $program (x work dirs | grep -w "$selected\$")
end
