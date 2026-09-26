package handler

import "fmt"

func handlePing(command Command) []byte {
	// No arguments: reply with a simple string
	if len(command.Args) == 0 {
		return []byte("+PONG\r\n")
	}

	// Only one argument allowed
	if len(command.Args) > 1 {
		return []byte("-ERR wrong number of arguments for 'ping' command\r\n")
	}

	// Return bulk string containing message
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(command.Args[0]), command.Args[0]))
}
