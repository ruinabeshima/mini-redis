# RESP Package

This package handles the parsing of incoming byte streams and returns the appropriate command.
The byte stream is in RESP2 format.

## Data Types

Data types are determined by the first byte of the stream.

- Simple strings (+)
- Simple errors (-)
- Integers (:)
- Bulk Strings ($)
- Arrays (\*)

## Architecture¡

- **`types.go`**: The main entry point of the package.
  - Defines constants for data and error types.
  - Contains `Parse()` which consumes the first byte of the stream and runs `switch` logic to distinguish between data types.
  - Calls specific `parser.go` functions corresponding to the data type.
  - The `Parse()` function returns a struct `Value`:
    ```go
    type Value struct {
        Type   byte   // '+', '-', ':', '$', '*'
        Str    string // Simple string, simple error, bulk string
        Int    int64
        Array  []Value
        IsNull bool // Null bulk strings, null array
    }
    ```

- **`parser.go`**: Contains functions to parse byte streams according to data type.
  - Functions included:
    - `parseSimpleString()`
    - `parseSimpleError()`
    - `parseInteger()`
    - `parseBulkString()`
    - `parseArray()`
  - Each parser function takes in a `bufio.Reader` as a parameter and returns the decoded value of the corresponding data type to put into the struct `Value`.
  - `parseArray()` recursively parses each child element by calling `Parse()`.

- **`helpers.go`**: Contains the helper function `readLine()`, which reads a line of data until a trailing newline and verifies that it ends with `\r\n`.

- **`writer.go`**: Defines a `Writer` struct with methods.
  - Function `NewWriter()` acts as a constructor; it wraps the given destination in a `bufio.Writer` and returns a ready to use pointer `*Writer`. It accepts any `io.Writer`, not just a network connection, so tests can write to an in-memory buffer and check the bytes without opening a socket.
  - A client can pipeline, sending several commands in one go before reading any replies. `server.go` therefore delays `Flush()` until `reader.Buffered()` reports 0,
  meaning every command already received has been handled and the client must now be waiting. A batch of pipelined commands then costs one write to the socket instead of one per command.
  - `server.go` calls `Flush()` only when `reader.Buffered()` returns 0, meaning no further pipelined commands are waiting and the client must be blocked on a reply. This batches a run of pipelined commands into a single write to the socket instead of one write per command.
  - Function `WriteValue()` takes in struct `Value` as input, obtains the data type and runs a switch case to write different types to the buffer.
  - The writer handles null values for arrays and bulk strings, and returns their respective null responses.
  - Function `writeArray()` recurses by calling `WriteValue()` for every element in its array.
  - Methods:
    - `writeSimpleString()`
    - `writeSimpleError()`
    - `writeInteger()`
    - `writeBulkString()`
    - `writeArray()`

  ```go
    type Writer struct {
      w *bufio.Writer
    }
  ```

- **`parser_test.go`**: Runs table-driven tests for the `resp` package (100% coverage reached).

### Architecture Diagram

```text
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

## Limits

The parser sets a limit for the length of bulk strings (512MB) and the number of array elements (1,024,000)
