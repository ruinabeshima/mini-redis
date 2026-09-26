package handler 

import "strings"

type Command struct {
	Name string 
	Args []string
}

func ParseCommand(parsedArray []string) Command {
	var command Command 

	// First element is the command 
	command.Name = parsedArray[0]

	// Other elements after are the commands
	if len(parsedArray) > 1 {
		command.Args = parsedArray[1:]
	}

	return command 
}