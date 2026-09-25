# Mini Redis

A Redis server written from scratch in Go, which listens on TCP `6379`.
As of now, every read still returns `+PONG\r\n` back, regardless of what the client
actually sends.

## Architecture

### Root

- [`main.go`](main.go) binds the listener — `net.Listen("tcp", ":6379")` — and
  accepts in a loop. A failed`Accept` is logged and the loop continues, so one bad connection can't take
  the server down.
- [`server.go`](server.go) runs one `handleConnection` goroutine per client.

### The read loop

There are two buffers, in `server.go`:

- `readBuffer` — a fixed 1 KiB scratch buffer, reused on every `Read`.
- `streamBuffer` — a growing buffer holding everything received but not yet
  parsed.

The inner loop parses _as many complete messages as the buffer holds_ and keeps
whatever is left over. Each successful parse returns how many bytes
it used, and `streamBuffer` is resliced past them, so a half-received command
simply stays in the buffer until the rest of it arrives.

## Packages

### `resp` — RESP2 parser

[`package resp`](resp/) reads **RESP2** (the wire format Redis uses) from raw
bytes and turns it into Go values. Dependency-free: stdlib only.

It is a package inside this module, not a module of its own. It was briefly
split into a separate repo and folded back in, so its import path is
`github.com/ruinabeshima/mini-redis/resp`.

#### How it works

`Parse` looks at the first byte to decide what the message is, hands it to the matching helper, and returns the value plus how many bytes it used.

```
 bytes off the connection
            │
            ▼
      Parse(data)  ── looks at data[0]
            │
   ┌────────┼──────────────────────────────┐
   │ '+'  simple string                    │
   │ '-'  error                            │
   │ ':'  integer                          │──► (Value, bytesRead, error)
   │ '$'  bulk string  (read by length)    │
   │ '*'  array ──┐                        │
   └──────────────┼────────────────────────┘
                  │  calls Parse on each element,
                  └─ adding up bytesRead as it goes
                        (nested arrays just work)
```

Every helper returns `bytesRead` so the caller knows where the next message starts in the buffer. If the bytes end mid-message, you get `ErrIncomplete` — read more from the socket and try again.

#### Files

| file | holds |
| --- | --- |
| [`resp/types.go`](resp/types.go) | the `Value` struct, the type-byte constants, `ErrIncomplete`, and `Parse` |
| [`resp/parser.go`](resp/parser.go) | one unexported helper per RESP2 type |
| [`resp/helpers.go`](resp/helpers.go) | CRLF / line-reading utilities shared by the helpers |

#### Exported surface

`Parse`, `Value`, and `ErrIncomplete`. The per-type helpers and the type-byte
constants are unexported — callers switch on `Value.Type` against the literal
bytes `'+' '-' ':' '$' '*'`.

#### The `Value` type

One struct holds any RESP2 message:

```go
type Value struct {
	Type   byte    // '+', '-', ':', '$', or '*'
	Str    string  // simple strings, errors, bulk strings
	Int    int     // integers
	Array  []Value // elements, for arrays
	IsNull bool    // $-1 or *-1
}
```

#### Supported types

Simple strings (`+`), errors (`-`), integers (`:`), bulk strings (`$`, including empty and null), and arrays (`*`, including nested and null).

Bulk strings are read by their declared byte length rather than scanned for a newline, so binary payloads pass through untouched.

#### Usage

No `go get` — the package ships with this module. Within the module, import it directly:

```go
import "github.com/ruinabeshima/mini-redis/resp"
```

The server consumes it in [`server.go`](server.go), which feeds it a growing
`[]byte` and reslices by the returned byte count (see [The read loop](#the-read-loop)).

## Getting Started

```sh
go build -o bin/mini-redis .
./bin/mini-redis
# Server started on PORT 6379
```

Port 6379 must be free. The module targets Go 1.27.1 (see [`go.mod`](go.mod)).
