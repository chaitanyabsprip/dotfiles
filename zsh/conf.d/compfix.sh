#!/bin/zsh

compgen () {
    local opts prefix suffix job OPTARG OPTIND ret=1
    local -a name res results jids
    local -A shortopts
    emulate -L sh
    setopt kshglob noshglob braceexpand nokshautoload
    shortopts=(a alias b builtin c command d directory e export f file g group j job k keyword u user v variable)
    while getopts "o:A:G:C:F:P:S:W:X:abcdefgjkuv" name
    do
        case $name in
            ([abcdefgjkuv]) OPTARG="${shortopts[$name]}"  ;&
            (A) case $OPTARG in
                    (alias) results+=("${(k)aliases[@]}")  ;;
                    (arrayvar) results+=("${(k@)parameters[(R)array*]}")  ;;
                    (binding) results+=("${(k)widgets[@]}")  ;;
                    (builtin) results+=("${(k)builtins[@]}" "${(k)dis_builtins[@]}")  ;;
                    (command) results+=("${(k)commands[@]}" "${(k)aliases[@]}" "${(k)builtins[@]}" "${(k)functions[@]}" "${(k)reswords[@]}")  ;;
                    (directory) setopt bareglobqual
                        results+=(${IPREFIX}${PREFIX}*${SUFFIX}${ISUFFIX}(N-/))
                        setopt nobareglobqual ;;
                    (disabled) results+=("${(k)dis_builtins[@]}")  ;;
                    (enabled) results+=("${(k)builtins[@]}")  ;;
                    (export) results+=("${(k)parameters[(R)*export*]}")  ;;
                    (file) setopt bareglobqual
                        results+=(${IPREFIX}${PREFIX}*${SUFFIX}${ISUFFIX}(N))
                        setopt nobareglobqual ;;
                    (function) results+=("${(k)functions[@]}")  ;;
                    (group) emulate zsh
                        _groups -U -O res
                        emulate sh
                        setopt kshglob noshglob braceexpand
                        results+=("${res[@]}")  ;;
                    (hostname) emulate zsh
                        _hosts -U -O res
                        emulate sh
                        setopt kshglob noshglob braceexpand
                        results+=("${res[@]}")  ;;
                    (job) results+=("${savejobtexts[@]%% *}")  ;;
                    (keyword) results+=("${(k)reswords[@]}")  ;;
                    (running) jids=("${(@k)savejobstates[(R)running*]}")
                        for job in "${jids[@]}"
                        do
                            results+=(${savejobtexts[$job]%% *})
                        done ;;
                    (stopped) jids=("${(@k)savejobstates[(R)suspended*]}")
                        for job in "${jids[@]}"
                        do
                            results+=(${savejobtexts[$job]%% *})
                        done ;;
                    (setopt | shopt) results+=("${(k)options[@]}")  ;;
                    (signal) results+=("SIG${^signals[@]}")  ;;
                    (user) results+=("${(k)userdirs[@]}")  ;;
                    (variable) results+=("${(k)parameters[@]}")  ;;
                    (helptopic)  ;;
                esac ;;
            (F) COMPREPLY=()
                local -a args
                args=("${words[0]}" "${@[-1]}" "${words[CURRENT-2]}")
                () {
                    typeset -h words
                    $OPTARG "${args[@]}"
                }
                results+=("${COMPREPLY[@]}")  ;;
            (G) setopt nullglob
                results+=(${~OPTARG})
                unsetopt nullglob ;;
            (W) results+=(${(Q)~=OPTARG})  ;;
            (C) results+=($(eval $OPTARG))  ;;
            (P) prefix="$OPTARG"  ;;
            (S) suffix="$OPTARG"  ;;
            (X) if [[ ${OPTARG[0]} = '!' ]]
                then
                    results=("${(M)results[@]:#${OPTARG#?}}")
                else
                    results=("${results[@]:#$OPTARG}")
                fi ;;
        esac
    done
    print -l -r -- "$prefix${^results[@]}$suffix"
}

_bash_complete () {
    local ret=1
    local -a suf matches
    local -x COMP_POINT COMP_CWORD
    local -a COMP_WORDS COMPREPLY BASH_VERSINFO
    local -x COMP_LINE="$words"
    local -A savejobstates savejobtexts
    (( COMP_POINT = 1 + ${#${(j. .)words[1,CURRENT-1]}} + $#QIPREFIX + $#IPREFIX + $#PREFIX ))
    (( COMP_CWORD = CURRENT - 1))
    COMP_WORDS=("${words[@]}")
    BASH_VERSINFO=(2 05b 0 1 release)
    savejobstates=(${(kv)jobstates})
    savejobtexts=(${(kv)jobtexts})
    [[ ${argv[${argv[(I)nospace]:-0}-1]} = -o ]] && suf=(-S '')
    matches=(${(f)"$(compgen $@ -- ${words[CURRENT]})"})
    if [[ -n $matches ]]
    then
        if [[ ${argv[${argv[(I)filenames]:-0}-1]} = -o ]]
        then
            compset -P '*/' && matches=(${matches##*/})
            compset -S '/*' && matches=(${matches%%/*})
            compadd -f "${suf[@]}" -a matches && ret=0
        else
            compadd "${suf[@]}" - "${(@)${(Q@)matches}:#*\ }" && ret=0
            compadd -S ' ' - ${${(M)${(Q)matches}:#*\ }% } && ret=0
        fi
    fi
    if (( ret ))
    then
        if [[ ${argv[${argv[(I)default]:-0}-1]} = -o ]]
        then
            _default "${suf[@]}" && ret=0
        elif [[ ${argv[${argv[(I)dirnames]:-0}-1]} = -o ]]
        then
            _directories "${suf[@]}" && ret=0
        fi
    fi
    return ret
}

complete () {
    emulate -L zsh
    local args void cmd print remove
    args=("$@")
    zparseopts -D -a void o: A: G: W: C: F: P: S: X: a b c d e f g j k u v p=print r=remove
    if [[ -n $print ]]
    then
        printf 'complete %2$s %1$s\n' "${(@kv)_comps[(R)_bash*]#* }"
    elif [[ -n $remove ]]
    then
        for cmd
        do
            unset "_comps[$cmd]"
        done
    else
        compdef _bash_complete\ ${(j. .)${(q)args[1,-1-$#]}} "$@"
    fi
}

# This file includes code from github.com/ohmyzsh/ohmyzsh, licensed under the MIT License.

# MIT License
#
# Copyright (c) 2009-2022 Robby Russell and contributors (https://github.com/ohmyzsh/ohmyzsh/contributors)
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
