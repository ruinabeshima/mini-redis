/*
	Entry point of package: parses incoming byte streams in RESP2 Format and returns command
	Includes simple strings, simple errors, integers, bulkStrings, array
	Parse() function is exported into main/server.go
*/

package resp

import (
	"bufio"
	"errors"
)

// Data types correspond to symbol of first byte
const (
	simpleString = '+'
	simpleError  = '-'
	integer      = ':'
	bulkString   = '$'
	array        = '*'
)

// Redis bulk string length limit (512MB)
const maxBulkLength = 512 * 1024 * 1024

// Redis array element limit (1024000 elements)
const maxArrayLength = 1024000

type Value struct {
	Type   byte   // '+', '-', ':', '$', '*'
	Str    string // Simple string, simple error, bulk string
	Int    int
	Array  []Value
	IsNull bool // Null bulk strings, null array
}

// Export error types
var ErrIncomplete = errors.New("incomplete RESP payload")
var ErrInvalidLength = errors.New("negative lengths less than -1 not permitted")
var ErrUnknownType = errors.New("unknown / invalid command")
var ErrTooLarge = errors.New("length exceeds limit")

func Parse(r *bufio.Reader) (Value, int, error) {

	// Read and consume the first byte (command type)
	typeByte, err := r.ReadByte()
	if err != nil {
		return Value{}, 0, err
	}

	switch typeByte {
	case simpleString:
		str, consumed, err := parseSimpleString(r)
		return Value{Type: simpleString, Str: str}, consumed, err

	case simpleError:
		errStr, consumed, err := parseSimpleError(r)
		return Value{Type: simpleError, Str: errStr}, consumed, err

	case integer:
		num, consumed, err := parseInteger(r)
		return Value{Type: integer, Int: num}, consumed, err

	case bulkString:
		bstr, isNull, consumed, err := parseBulkString(r)
		return Value{Type: bulkString, IsNull: isNull, Str: bstr}, consumed, err

	case array:
		arr, consumed, err := parseArray(r)
		return arr, consumed, err

	default:
		return Value{}, 0, ErrUnknownType
	}
}
