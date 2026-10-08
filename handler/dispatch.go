/*
	Dispatcher that uses a hashmap for direct lookup
	Receives parsed command and arguments from ParseCommand(), and executes said command
*/

package handler

import "fmt"

// Hashmap with string keys and function values
type handlerFunc func(args []string) []byte

var operations = map[string]handlerFunc{
	"PING": handlePing,
	"ECHO": handleEcho,
	"SET":  handleSet,
}

func ExecuteCommand(command Command) []byte {
	fn, exists := operations[command.Name]

	if !exists {
		return []byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command.Name))
	}

	return fn(command.Args)
}
