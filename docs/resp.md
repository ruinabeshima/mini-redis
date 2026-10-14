# RESP Package 

This package handles the parsing of incoming byte streams and returns the appropriate command. 
The byte stream is in RESP2 format. 

## Data Types 
Data types are determined by the first byte of the stream. 
- Simple strings (+)
- Simple errors (-)
- Integers (:)
- Bulk Strings ($)
- Arrays (*)


## Architecture

* **`types.go`**: The main entry point of the package.
  * Defines constants for data and error types.
  * Contains `Parse()` which consumes the first byte of the stream and runs `switch` logic to distinguish between data types.
  * Calls specific `parser.go` functions corresponding to the data type.
  * The `Parse()` function returns a struct `Value`:
    ```go
    type Value struct {
        Type   byte   // '+', '-', ':', '$', '*'
        Str    string // Simple string, simple error, bulk string
        Int    int
        Array  []Value
        IsNull bool // Null bulk strings, null array
    }
    ```

* **`parser.go`**: Contains functions to parse byte streams according to data type.
  * Functions included:
    * `parseSimpleString()`
    * `parseSimpleError()`
    * `parseInteger()`
    * `parseBulkString()`
    * `parseArray()`
  * Each parser function takes in a `bufio.Reader` as a parameter and returns the decoded value of the corresponding data type to put into the struct `Value`.
  * `parseArray()` recursively parses each child element by calling `Parse()`.

* **`helpers.go`**: Contains the helper function `readLine()`, which reads a line of data until a trailing newline and verifies that it ends with `\r\n`.

* **`parser_test.go`**: Runs table-driven tests for the `resp` package (100% coverage reached).

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


