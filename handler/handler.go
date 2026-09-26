package handler 

import (
	"errors"
	"strings"
)

type Command struct {
	Name string 
	Args []string
}

func ParseCommand(parsedArray []string) Command, error {
	var command Command 

	if len(parsedArray) == 0 {
		return Command{}, errors.New("no command provided")
	}

	// First element is the command 
	command.Name = strings.ToUpper(parsedArray[0])

	// Other elements after are the commands
	if len(parsedArray) > 1 {
		command.Args = parsedArray[1:]
	}

	return command 
}