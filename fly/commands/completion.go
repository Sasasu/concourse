package commands

import (
	"fmt"
	"reflect"
)

type CompletionCommand struct {
	Shell string `long:"shell" required:"true" choice:"bash" choice:"zsh" choice:"fish"` // add more choices later
}

// credits:
// https://godoc.org/github.com/jessevdk/go-flags#hdr-Completion
// https://github.com/concourse/concourse/issues/1309#issuecomment-452893900
const bashCompletionSnippet = `_fly_compl() {
	args=("${COMP_WORDS[@]:1:$COMP_CWORD}")
	local IFS=$'\n'
	COMPREPLY=($(GO_FLAGS_COMPLETION=1 ${COMP_WORDS[0]} "${args[@]}"))
	return 0
}
complete -F _fly_compl fly
`

func fishCompletionSnippetHelper(snippet string, prefix string, commandType reflect.Type) string {
	for i := 0; i < commandType.NumField(); i++ {
		var tags = commandType.Field(i).Tag

		var command = tags.Get("command")
		var long, short = tags.Get("long"), tags.Get("short")
		var description = tags.Get("description")

		var template = "complete --no-files -c fly"

		if prefix != "" {
			template += fmt.Sprintf(" -n '__fish_seen_subcommand_from %s'", prefix)
		}

		if command != "" {
			snippet += fishCompletionSnippetHelper(template, prefix+" "+command, commandType.Field(i).Type)

			template += fmt.Sprintf(" -n '__fish_use_subcommand' -a %s", command)
		}

		if description != "" {
			template += fmt.Sprintf(` -d "%s"`, description)
		}

		if long != "" {
			template += fmt.Sprintf(" --l %s", long)
		}

		if short != "" {
			template += fmt.Sprintf(" -s %s", short)
		}

		snippet += template + "\n"
	}

	return snippet
}

var fishCompletionSnippet = fishCompletionSnippetHelper("", "", reflect.TypeOf(Fly))

// initial implemenation just using bashcompinit
const zshCompletionSnippet = `autoload -Uz compinit && compinit
autoload -Uz bashcompinit && bashcompinit
` + bashCompletionSnippet

func (command *CompletionCommand) Execute([]string) error {
	switch command.Shell {
	case "bash":
		_, err := fmt.Print(bashCompletionSnippet)
		return err
	case "zsh":
		_, err := fmt.Print(zshCompletionSnippet)
		return err
	case "fish":
		_, err := fmt.Print(fishCompletionSnippet)
		return err
	default:
		// this should be unreachable
		return fmt.Errorf("unknown shell %s", command.Shell)
	}
}
