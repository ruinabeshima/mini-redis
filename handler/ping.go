package handler

import "github.com/ruinabeshima/mini-redis/resp"

func (h *Handler) handlePing(args []string) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.SimpleString, Str: "PONG"}
	}
	if len(args) > 1 {
		return resp.Value{Type: resp.SimpleError, Str: "ERR wrong number of arguments for 'ping' command"}
	}

	// Bulk string with argument
	return resp.Value{Type: resp.BulkString, Str: args[0]}
}

func (h *Handler) handleEcho(args []string) resp.Value {
	if len(args) == 0 || len(args) > 1 {
		return resp.Value{Type: resp.SimpleError, Str: "ERR wrong number of arguments for 'echo' command"}
	}

	return resp.Value{Type: resp.BulkString, Str: args[0]}
}

func (h *Handler) handleSet(args []string) resp.Value {
	if len(args) != 2 {
		return resp.Value{Type: resp.SimpleError, Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0]
	value := args[1]
	h.store.Set(key, value)
	return resp.Value{Type: resp.SimpleString, Str: "OK"}
}

func (h *Handler) handleGet(args []string) resp.Value {
	if len(args) != 1 {
		return resp.Value{Type: resp.SimpleError, Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0]
	val, ok := h.store.Get(key)
	if !ok {
		return resp.Value{Type: resp.BulkString, IsNull: true} // Null bulk string
	}

	return resp.Value{Type: resp.BulkString, Str: val}
}

func (h *Handler) handleDel(args []string) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.SimpleError, Str: "ERR wrong number of arguments for 'del' command"}
	}

	count := 0
	for _, key := range args {
		if h.store.Del(key) {
			count++
		}
	}

	return resp.Value{Type: resp.Integer, Int: int64(count)}
}

func (h *Handler) handleExists(args []string) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: resp.SimpleError, Str: "ERR wrong number of arguments for 'exists' command"}
	}

	count := 0
	for _, key := range args {
		if h.store.Exists(key) {
			count++
		}
	}

	return resp.Value{Type: resp.Integer, Int: int64(count)}
}
