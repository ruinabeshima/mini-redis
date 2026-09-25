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

- **[`resp`](resp/README.md)** — parses RESP2 from raw bytes into Go values:
  simple strings, simple errors, integers, bulk strings, and arrays, including
  nested and null forms. Every parse returns the number of bytes it consumed, so
  the caller knows where the next message starts.
  [Full documentation →](resp/README.md)

## Getting Started

```sh
go build -o bin/mini-redis .
./bin/mini-redis
# Server started on PORT 6379
```

Port 6379 must be free.
