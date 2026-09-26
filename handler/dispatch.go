/*
	Dispatcher that uses a hashmap for direct lookup
	Receives parsed command and arguments from ParseCommand(), and executes said command
*/

package handler

import "errors"

// Hashmap with string keys and function values
type handlerFunc func(args []string) []byte

var operations = map[string]handlerFunc{
	"PING": handlePing,
}

func ExecuteCommand(command Command) []byte {
	fn, exists = operations[command.Name]

	if !exists {
		return []byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd.Name))
	}

	return fn(command.Args)
}
