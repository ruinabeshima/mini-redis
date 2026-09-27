package handler

import "fmt"

// Tests server availability (health check)
func handlePing(args []string) []byte {
	if len(args) == 0 {
		return []byte("+PONG\r\n")
	}
	if len(args) > 1 {
		return []byte("-ERR wrong number of arguments for 'ping' command\r\n")
	}

	// Bulk string with argument
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0]), args[0]))
}

// Tests data transmission with specific message
func handleEcho(args []string) []byte {
	if len(args) == 0 || len(args) > 1 {
		return []byte("-ERR wrong number of arguments for 'echo' command\r\n")
	}

	// Bulk string with argument
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0]), args[0]))
}
