package commands

import (
	"fmt"
	"strings"

	"github.com/AbGuthrie/goquery/config"
	prompt "github.com/c-bata/go-prompt"
)

func alias(cmdline string) error {
	args := strings.Split(cmdline, " ")
	if len(args) == 1 {
		// No args provided, print current state of aliases
		aliases := config.GetConfig().Aliases
		if len(aliases) == 0 {
			fmt.Printf("No aliases set\n")
			return nil
		}
		fmt.Printf("Available aliases:\n\n")
		for _, alias := range aliases {
			fmt.Printf("Name: %s\nCommand: %s\n\n", alias.Name, alias.Command)
			return nil
		}
	}

	// Otherwise create a new alias
	args = args[1:]
	name := args[0]
	command := ""
	if len(args) > 1 {
		command = args[1]
	}

	// Create the command and store in state
	err := config.AddAlias(name, command)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error creating alias: %s\n", err))
	}

	fmt.Printf("Created new alias '%s' with command: %s\n", name, command)
	return nil
}

func aliasHelp() string {
	return "Create a new alias or call with no arguments to list current aliases. " +
		"The format for creating an alias is as follows: ALIAS_NAME .example arg1 $# arg3"
}

func aliasSuggest(cmdline string) []prompt.Suggest {
	suggestions := []prompt.Suggest{}
	for _, alias := range config.GetConfig().Aliases {
		suggestions = append(suggestions, prompt.Suggest{
			Text:        alias.Name,
			Description: alias.Command,
		})
	}
	return suggestions
}

// FindAlias searches the list of named aliases and returns the Alias struct if found
func FindAlias(command string) (config.Alias, bool) {
	aliases := config.GetConfig().Aliases
	for _, alias := range aliases {
		if command == alias.Name {
			return alias, true
		}
	}
	return config.Alias{}, false
}

// InterpolateArguments fills in an alias' placeholders ($#) with provided arguments
// TODO add alias_test.go unit tests
func InterpolateArguments(rawLine string, command string) (string, error) {
	inputParts := strings.Split(rawLine, " ")
	args := inputParts[1:]

	// TODO this should support escaping and ignoring the
	// placeholder pattern ie \$#
	placeholderParts := strings.Split(command, "$#")

	// Assert arguments provided and placeholders align
	if len(args) != len(placeholderParts)-1 {
		return "", fmt.Errorf("Argument mismatch, alias expects %d args", len(placeholderParts)-1)
	}

	// If no placeholders in query, return as is
	if len(placeholderParts)-1 == 0 {
		return command, nil
	}

	realizedCommand := ""
	for i, arg := range args {
		realizedCommand += placeholderParts[i]
		realizedCommand += arg
	}

	return realizedCommand, nil
}