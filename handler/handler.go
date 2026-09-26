package handler

import (
	"errors"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func ParseCommand(parsedArray []string) (Command, error) {
	var command Command

	if len(parsedArray) == 0 {
		return Command{}, errors.New("no command provided")
	}

	// First element: command name
	command.Name = strings.ToUpper(parsedArray[0])

	// Other elements: command arguments
	if len(parsedArray) > 1 {
		command.Args = parsedArray[1:]
	}

	return command, nil
}
