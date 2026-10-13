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

- `types.go`: The main entry point of the package.  
    - Defines constants for data and error types. 
    - Contains `Parse()` which consumes the first byte of the stream and runs `switch` logic to distinguish between data types. 
    - Calls specific `parser.go` functions corresponding to data type. 
    - The `Parse()` function returns a struct `Value`: 
    ``` 
        type Value struct {
            Type   byte   // '+', '-', ':', '$', '*'
            Str    string // Simple string, simple error, bulk string
            Int    int
            Array  []Value
            IsNull bool // Null bulk strings, null array
        }
    ```
- `parser.go`: Contains functions to parse byte stream according to data type. 
    - Functions included: 
        - `parseSimpleString()`
        - `parseSimpleError()`
        - `parseInteger()`
        - `parseBulkString()`
        - `parseArray()`
    - Each parser function takes in a Bufio reader as a parameter, and returns the decoded value of the corresponding data type to put into struct `Value`. 
    - `parseArray()` recursively parses each child element by calling `Parse()`. 
- `helpers.go`: Contains helper function `readLine()` which reads line of data until trailing \n, and verifies that it ends with \r\n. 
- `parser_test.go`: Runs table-driven tests for the `resp` package. 100% coverage reached. 

### Architecture Diagram 
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

## Error types 
- `ErrInvalidLength`: Negative lengths (for bulk strings and arrays) less than -1 are not permitted
- `ErrUnknownType`: Unknown / invalid command 
- `ErrTooLarge`: Length exceeds limit (512MB for bulk strings, 1024000 array elements)
- `ErrNoCRLF`: No trailing CRLF (\r\n)

## Limits 
The parser sets a limit for the length of bulk strings (512MB) and the number of array elements (1,024,000)


