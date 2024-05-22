package main

// Print out shell completions that can be sourced from shell init files

import (
	"fmt"
)

const zsh_comp = `
autoload -Uz compinit
compinit -u

_toru() {
    local -a completions
    args=("${words[@]:1}")
    local IFS=$'\n'
    completions=($(GO_FLAGS_COMPLETION=1 ${words[1]} "${args[@]}"))
    compadd -a completions
}

compdef _toru toru
    `

const bash_comp = `
_completion_toru() {
    # All arguments except the first one
    args=("${COMP_WORDS[@]:1:$COMP_CWORD}")

    # Only split on newlines
    local IFS=$'\n'

    # Call completion (note that the first element of COMP_WORDS is
    # the executable itself)
    COMPREPLY=($(GO_FLAGS_COMPLETION=verbose ${COMP_WORDS[0]} "${args[@]}"))
    return 0
}

complete -F _completion_toru toru
`

func Completers() error {
	if completions.Zsh {
		fmt.Println(zsh_comp)
	} else if completions.Bash {
		fmt.Println(bash_comp)
	}

	return fmt.Errorf("Must choose a shell [bash|zsh] to output completions for")
}
