/*
	Dispatcher that uses a hashmap for direct lookup
	Receives parsed command and arguments from ParseCommand(), and executes said command
*/

package handler

import (
	"fmt"
	"github.com/ruinabeshima/mini-redis/resp"
	"github.com/ruinabeshima/mini-redis/store"
)

// Hashmap with string keys and function values
type handlerFunc func(args []string) resp.Value

type Handler struct {
	store      *store.Store
	operations map[string]handlerFunc
}

func New(s *store.Store) *Handler {
	h := &Handler{store: s}
	h.operations = map[string]handlerFunc{
		"PING": h.handlePing,
		"ECHO": h.handleEcho,
		"SET":  h.handleSet,
		"GET":  h.handleGet,
	}
	return h
}

func (h *Handler) ExecuteCommand(command Command) resp.Value {
	fn, exists := h.operations[command.Name]

	if !exists {
		return resp.Value{Type: resp.SimpleError, Str: fmt.Sprintf("ERR unknown command '%s'", command.Name)}
	}

	return fn(command.Args)
}
