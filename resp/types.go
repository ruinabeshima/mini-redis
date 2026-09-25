/* 
	Entry point of package: parses incoming byte streams in RESP2 Format and returns command
	Includes simple strings, simple errors, integers, bulkStrings, array 
	Parse() function is exported into main/connection.go
*/

package resp

import "errors"

// Data types correspond to symbol of first byte
const (
	simpleString = '+'
	simpleError  = '-'
	integer      = ':'
	bulkString   = '$'
	array        = '*'
)

type Value struct {
	Type   byte   // '+', '-', ':', '$', '*'
	Str    string // Simple string, simple error, bulk string
	Int    int
	Array  []Value
	IsNull bool // Null bulk strings, null array
}

func Parse(data []byte) (Value, int, error) {
	if len(data) == 0 {
		return Value{}, 0, errors.New("empty payload")
	}

	switch data[0] {
	case simpleString:
		str, consumed, err := parseSimpleString(data)
		return Value{Type: simpleString, Str: str}, consumed, err

	case simpleError:
		errStr, consumed, err := parseSimpleError(data)
		return Value{Type: simpleError, Str: errStr}, consumed, err

	case integer:
		num, consumed, err := parseInteger(data)
		return Value{Type: integer, Int: num}, consumed, err

	case bulkString:
		bstr, isNull, consumed, err := parseBulkString(data)
		return Value{Type: bulkString, IsNull: isNull, Str: bstr}, consumed, err

	case array:
		arr, consumed, err := parseArray(data)
		return arr, consumed, err

	default:
		return Value{}, 0, errors.New("unknown / invalid command")
	}
}
