# RESP2 Parser

`package resp` — reads **RESP2** (the wire format Redis uses) from raw bytes and
turns it into Go values. Dependency-free: stdlib only.

This is a package inside the [`github.com/ruinabeshima/mini-redis`](https://github.com/ruinabeshima/mini-redis)
module, not a module of its own. It was briefly split into a separate repo and
folded back in, so its import path is now `github.com/ruinabeshima/mini-redis/resp`.

## How it works

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

## Files

| file | holds |
| --- | --- |
| `types.go` | the `Value` struct, the type-byte constants, `ErrIncomplete`, and `Parse` |
| `parser.go` | one unexported helper per RESP2 type |
| `helpers.go` | CRLF / line-reading utilities shared by the helpers |

## Exported surface

`Parse`, `Value`, and `ErrIncomplete`. The per-type helpers and the type-byte
constants are unexported — callers switch on `Value.Type` against the literal
bytes `'+' '-' ':' '$' '*'`.

## The `Value` type

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

## Supported types

Simple strings (`+`), errors (`-`), integers (`:`), bulk strings (`$`, including empty and null), and arrays (`*`, including nested and null).

Bulk strings are read by their declared byte length rather than scanned for a newline, so binary payloads pass through untouched.

## Usage

No `go get` — the package ships with this module. Within the module, import it directly:

```go
import "github.com/ruinabeshima/mini-redis/resp"
```

The server consumes it in [`server.go`](../server.go), which feeds it a growing
`[]byte` and reslices by the returned byte count.

Module targets Go 1.27.1 (see [`go.mod`](../go.mod)).
