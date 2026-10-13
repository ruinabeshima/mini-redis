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
