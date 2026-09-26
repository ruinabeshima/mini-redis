package handler

import "fmt"

func handlePing(args []string) []byte {
	// No arguments: reply with a simple string
	if len(args) == 0 {
		return []byte("+PONG\r\n")
	}

	// Only one argument allowed
	if len(args) > 1 {
		return []byte("-ERR wrong number of arguments for 'ping' command\r\n")
	}

	// Return bulk string containing message
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0]), args[0]))
}
