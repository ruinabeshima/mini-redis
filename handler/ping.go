package handler

import (
	"fmt"
	"github.com/ruinabeshima/mini-redis/store"
)

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

// Adds to key-value store
func handleSet(args []string) []byte {
	if len(args) != 2 {
		return []byte("-ERR wrong number of arguments for 'set' command\r\n")
	}

	key := args[0]
	value := args[1]
	store.KVStore[key] = value
	return []byte("+OK\r\n")
}

// Retrieves from key-value store
func handleGet(args []string) []byte {
	if len(args) != 1 {
		return []byte("-ERR wrong number of arguments for 'set' command\r\n")
	}

	key := args[0]
	val, ok := store.KVStore[key]

	if ok != true {
		return []byte("$-1\r\n") // Return null bulk string
	}

	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val))
}
