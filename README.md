# Mini Redis

A Redis server written from scratch in Go, which listens on TCP `6379`.
It parses incoming **RESP2** commands and currently supports `PING` and `ECHO`.

## Architecture

### Root

- [`main.go`](main.go) binds the listener — `net.Listen("tcp", ":6379")` — and
  accepts in a loop. A failed `Accept` is logged and the loop continues, so one bad connection can't take
  the server down.
- [`server.go`](server.go) runs one `handleConnection` goroutine per client.

### The read loop

Each connection is wrapped in a single `bufio.Reader`, and `handleConnection`
loops over three steps:

1. **Parse** — `resp.Parse(reader)` reads exactly one complete RESP value. If the
   bytes haven't all arrived yet, it blocks on the socket until they do.
2. **Interpret** — `handler.ParseCommand` turns the value into a `Command`
   (name + args).
3. **Execute** — `handler.ExecuteCommand` looks up the command and writes the
   reply back to the connection.

Commands are handled one at a time, in order. Pipelined commands that arrive in
one TCP packet stay in the reader's buffer and are picked up on the next loop.

Errors are handled differently depending on where they happen:

| where | what happens |
| --- | --- |
| `Parse` returns `io.EOF` | the client disconnected, so the connection is closed without an error reply |
| any other `Parse` error | reply `-ERR <msg>`, then close. The stream can't be trusted after bad input |
| `ParseCommand` error | reply `-ERR <msg>` and keep the connection open |
| a write fails | log it and close |

## Packages

### `resp` — RESP2 parser

[`package resp`](resp/) reads **RESP2** (the wire format Redis uses) from a
`*bufio.Reader` and turns it into Go values. Dependency-free: stdlib only.

It is a package inside this module, not a module of its own. It was briefly
split into a separate repo and folded back in, so its import path is
`github.com/ruinabeshima/mini-redis/resp`.

#### How it works

`Parse` reads and consumes the first byte to work out what type the message is,
then hands the reader to the matching helper. The reader moves forward as it
reads, so no byte counts are passed around. When `Parse` returns, the reader is
positioned at the start of the next message.

```
     *bufio.Reader (wraps net.Conn)
            │
            ▼
      Parse(r)  ── r.ReadByte()
            │
   ┌────────┼──────────────────────────────┐
   │ '+'  simple string                    │
   │ '-'  error                            │
   │ ':'  integer                          │──► (Value, error)
   │ '$'  bulk string  (read by length)    │
   │ '*'  array ──┐                        │
   └──────────────┼────────────────────────┘
                  │  calls Parse(r) once per element
                  └─ (nested arrays just work)
```

#### Files

| file | holds |
| --- | --- |
| [`resp/types.go`](resp/types.go) | the `Value` struct, the type-byte constants, the size limits, the exported errors, and `Parse` |
| [`resp/parser.go`](resp/parser.go) | one unexported helper per RESP2 type |
| [`resp/helpers.go`](resp/helpers.go) | `readLine`: reads up to `\n` and checks that the line ends in `\r\n` |
| [`resp/parser_test.go`](resp/parser_test.go) | table-driven tests for `Parse` |

#### Exported surface

`Parse`, `Value`, and these errors:

| error | when |
| --- | --- |
| `ErrUnknownType` | the first byte isn't one of `+ - : $ *` |
| `ErrNoCRLF` | a line ends in a bare `\n`, or a bulk string's payload isn't followed by `\r\n` |
| `ErrInvalidLength` | a bulk string or array length is below `-1` |
| `ErrTooLarge` | a bulk string is over 512 MB, or an array has more than 1,024,000 elements |

Errors from the reader (for example `io.EOF`) and from `strconv` pass through
unchanged. The per-type helpers and the type-byte constants are unexported, so
callers switch on `Value.Type` against the literal bytes `'+' '-' ':' '$' '*'`.

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

Bulk strings are read by their declared byte length with `io.ReadFull` rather
than scanned for a newline, so binary payloads pass through untouched.

#### Usage

No `go get` — the package ships with this module. Within the module, import it directly:

```go
import "github.com/ruinabeshima/mini-redis/resp"

reader := bufio.NewReader(conn)
val, err := resp.Parse(reader)
```

### `handler` — commands

[`package handler`](handler/) turns a parsed `resp.Value` into a reply.

| file | holds |
| --- | --- |
| [`handler/command.go`](handler/command.go) | `Command` and `ParseCommand`. It requires a non-empty array of bulk strings, and upper-cases the first element as the command name |
| [`handler/dispatch.go`](handler/dispatch.go) | `ExecuteCommand`, a map lookup from name to handler func. An unknown name gets `-ERR unknown command '<name>'` |
| [`handler/ping.go`](handler/ping.go) | `PING [message]` and `ECHO message` |

To add a command, write a `func(args []string) []byte` that returns the raw RESP
reply, and register it in the `operations` map.

## Getting Started

```sh
go build -o bin/mini-redis .
./bin/mini-redis
# Server started on PORT 6379
```

Port 6379 must be free. The module targets Go 1.27.1 (see [`go.mod`](go.mod)).

Try it with `redis-cli`:

```sh
redis-cli PING          # PONG
redis-cli ECHO hello    # "hello"
```

Run the tests:

```sh
go test ./...
```

## Architectural Decisions

### Bufio

The parser used to take a `[]byte` and return how many bytes it consumed. The
server kept a growing `streamBuffer`, resliced it past each parsed message, and
relied on `ErrIncomplete` to know when to read more from the socket.

It now takes a `*bufio.Reader` instead:

- **No manual buffer management.** The reader keeps track of position, so there's no
  `bytesRead` to add up through nested arrays, and no reslicing in the server.
- **Partial messages block instead of failing.** A half-received command simply
  waits for more bytes from the socket, so `ErrIncomplete` and the
  retry loop are gone.
- **Bulk strings read exactly what they declare.** `io.ReadFull` reads the payload,
  then a 2-byte read checks for the trailing `\r\n`.
- **Simpler tests.** Each test case wraps its input in a `bufio.Reader` and
  checks the parsed value plus `expectedRemaining`, the bytes left unread. This
  confirms that `Parse` consumes exactly one message.

The trade-off is that `Parse` is now tied to a stream. To parse a byte slice you
have to wrap it first, with `bufio.NewReader(bytes.NewReader(b))`.
