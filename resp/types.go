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
var ErrInvalidLength = errors.New("negative lengths less than -1 not permitted")
var ErrUnknownType = errors.New("unknown / invalid command")
var ErrTooLarge = errors.New("length exceeds limit")
var ErrNoCRLF = errors.New("line must end with CRLF")

func Parse(r *bufio.Reader) (Value, error) {

	// Read and consume the first byte (command type)
	typeByte, err := r.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch typeByte {
	case simpleString:
		str, err := parseSimpleString(r)
		return Value{Type: simpleString, Str: str}, err

	case simpleError:
		errStr, err := parseSimpleError(r)
		return Value{Type: simpleError, Str: errStr}, err

	case integer:
		num, err := parseInteger(r)
		return Value{Type: integer, Int: num}, err

	case bulkString:
		bstr, isNull, err := parseBulkString(r)
		return Value{Type: bulkString, IsNull: isNull, Str: bstr}, err

	case array:
		arr, err := parseArray(r)
		return arr, err

	default:
		return Value{}, ErrUnknownType
	}
}
