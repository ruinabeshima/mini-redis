/*
	Parses the command from an array of strings into a Go structure with name and arguments
	First element of array is the name
	All remaining elements are the arguments
*/

package handler

import (
	"errors"
	"github.com/ruinabeshima/mini-redis/resp"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func ParseCommand(val resp.Value) (Command, error) {

	// Command must be an array of bulk strings
	if val.Type != '*' || val.IsNull || len(val.Array) == 0 {
		return Command{}, errors.New("expected non-empty array of bulk strings")
	}

	args := make([]string, len(val.Array))
	for index, element := range val.Array {
		if element.Type != '$' || element.IsNull {
			return Command{}, errors.New("command arguments must be bulk string")
		}

		args[index] = element.Str
	}

	var command Command

	// First element: command name
	command.Name = strings.ToUpper(args[0])

	// Other elements: command arguments
	if len(val.Array) > 1 {
		command.Args = args[1:]
	}

	return command, nil
}
