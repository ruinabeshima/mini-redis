/*
	Table driven tests for Parse() function
	100% coverage reached
*/

package resp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {

	// Declare and initialise test cases
	tests := []struct {
		name              string
		input             []byte
		expectedVal       Value
		expectedRemaining string // Bytes left unread in the reader after parsing
		expectedError     error
	}{
		// Invalid input
		{
			name:          "empty payload",
			input:         []byte(""),
			expectedError: io.EOF,
		},
		{
			name:          "unknown type prefix",
			input:         []byte("!foo\r\n"),
			expectedError: ErrUnknownType,
		},
		{
			name:          "array element with unknown type prefix",
			input:         []byte("*1\r\n!foo\r\n"),
			expectedError: ErrUnknownType,
		},

		// Simple strings
		{
			name:        "simple string",
			input:       []byte("+OK\r\n"),
			expectedVal: Value{Type: simpleString, Str: "OK"},
		},
		{
			name:        "empty simple string",
			input:       []byte("+\r\n"),
			expectedVal: Value{Type: simpleString, Str: ""},
		},
		{
			name:              "simple string with trailing data only consumes first value",
			input:             []byte("+OK\r\n+PONG\r\n"),
			expectedVal:       Value{Type: simpleString, Str: "OK"},
			expectedRemaining: "+PONG\r\n",
		},
		{
			name:          "simple string missing CRLF",
			input:         []byte("+OK"),
			expectedError: io.EOF,
		},
		{
			name:          "simple string ending in bare LF",
			input:         []byte("+OK\n"),
			expectedError: ErrNoCRLF,
		},
		{
			name:          "simple string missing LF",
			input:         []byte("+OK\r"),
			expectedError: io.EOF,
		},

		// Simple errors
		{
			name:        "simple error",
			input:       []byte("-ERR unknown command\r\n"),
			expectedVal: Value{Type: simpleError, Str: "ERR unknown command"},
		},
		{
			name:          "simple error missing CRLF",
			input:         []byte("-ERR"),
			expectedError: io.EOF,
		},

		// Integers
		{
			name:        "positive integer",
			input:       []byte(":1000\r\n"),
			expectedVal: Value{Type: integer, Int: 1000},
		},
		{
			name:        "negative integer",
			input:       []byte(":-42\r\n"),
			expectedVal: Value{Type: integer, Int: -42},
		},
		{
			name:        "explicit plus sign integer",
			input:       []byte(":+5\r\n"),
			expectedVal: Value{Type: integer, Int: 5},
		},
		{
			name:        "zero integer",
			input:       []byte(":0\r\n"),
			expectedVal: Value{Type: integer, Int: 0},
		},
		{
			name:          "integer missing CRLF",
			input:         []byte(":12"),
			expectedError: io.EOF,
		},
		{
			name:          "non-numeric integer",
			input:         []byte(":abc\r\n"),
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "empty integer",
			input:         []byte(":\r\n"),
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "integer overflow",
			input:         []byte(":99999999999999999999\r\n"),
			expectedError: strconv.ErrRange,
		},

		// Bulk strings
		{
			name:        "bulk string",
			input:       []byte("$5\r\nhello\r\n"),
			expectedVal: Value{Type: bulkString, Str: "hello"},
		},
		{
			name:        "empty bulk string",
			input:       []byte("$0\r\n\r\n"),
			expectedVal: Value{Type: bulkString, Str: ""},
		},
		{
			name:        "bulk string containing CRLF is binary safe",
			input:       []byte("$7\r\nfoo\r\nba\r\n"),
			expectedVal: Value{Type: bulkString, Str: "foo\r\nba"},
		},
		{
			name:        "bulk string with multi-digit length",
			input:       []byte("$11\r\nhello world\r\n"),
			expectedVal: Value{Type: bulkString, Str: "hello world"},
		},
		{
			name:        "null bulk string",
			input:       []byte("$-1\r\n"),
			expectedVal: Value{Type: bulkString, IsNull: true},
		},
		{
			name:          "bulk string missing length CRLF",
			input:         []byte("$5"),
			expectedError: io.EOF,
		},
		{
			name:          "bulk string payload shorter than length",
			input:         []byte("$5\r\nhel"),
			expectedError: io.ErrUnexpectedEOF,
		},
		{
			name:          "bulk string missing trailing CRLF",
			input:         []byte("$5\r\nhello"),
			expectedError: io.EOF,
		},
		{
			name:          "bulk string trailing CR without LF",
			input:         []byte("$5\r\nhello\r"),
			expectedError: io.ErrUnexpectedEOF,
		},
		{
			name:          "bulk string wrong trailing bytes",
			input:         []byte("$5\r\nhelloXY"),
			expectedError: ErrNoCRLF,
		},
		{
			name:          "bulk string payload longer than length",
			input:         []byte("$3\r\nhello\r\n"),
			expectedError: ErrNoCRLF,
		},
		{
			name:          "bulk string length of -2",
			input:         []byte("$-2\r\n"),
			expectedError: ErrInvalidLength,
		},
		{
			name:          "bulk string large negative length",
			input:         []byte("$-100\r\nhello\r\n"),
			expectedError: ErrInvalidLength,
		},
		{
			name:          "bulk string length exceeds limit",
			input:         []byte("$9223372036854775807\r\n"),
			expectedError: ErrTooLarge,
		},
		{
			name:          "non-numeric bulk string length",
			input:         []byte("$abc\r\nhello\r\n"),
			expectedError: strconv.ErrSyntax,
		},

		// Arrays
		{
			name:        "empty array",
			input:       []byte("*0\r\n"),
			expectedVal: Value{Type: array, Array: []Value{}},
		},
		{
			name:        "null array",
			input:       []byte("*-1\r\n"),
			expectedVal: Value{Type: array, IsNull: true},
		},
		{
			name:  "array of bulk strings",
			input: []byte("*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n"),
			expectedVal: Value{Type: array, Array: []Value{
				{Type: bulkString, Str: "ECHO"},
				{Type: bulkString, Str: "hello"},
			}},
		},
		{
			name:  "array of mixed types",
			input: []byte("*3\r\n:1\r\n+OK\r\n$-1\r\n"),
			expectedVal: Value{Type: array, Array: []Value{
				{Type: integer, Int: 1},
				{Type: simpleString, Str: "OK"},
				{Type: bulkString, IsNull: true},
			}},
		},
		{
			name:  "nested array",
			input: []byte("*2\r\n*1\r\n:1\r\n+OK\r\n"),
			expectedVal: Value{Type: array, Array: []Value{
				{Type: array, Array: []Value{{Type: integer, Int: 1}}},
				{Type: simpleString, Str: "OK"},
			}},
		},
		{
			name:  "array with trailing data only consumes first value",
			input: []byte("*1\r\n+PING\r\n*1\r\n+PING\r\n"),
			expectedVal: Value{Type: array, Array: []Value{
				{Type: simpleString, Str: "PING"},
			}},
			expectedRemaining: "*1\r\n+PING\r\n",
		},
		{
			name:          "array missing length CRLF",
			input:         []byte("*2"),
			expectedError: io.EOF,
		},
		{
			name:          "array with incomplete element",
			input:         []byte("*2\r\n$4\r\nECHO\r\n$5\r\nhel"),
			expectedError: io.ErrUnexpectedEOF,
		},
		{
			name:          "array missing final element",
			input:         []byte("*2\r\n+OK\r\n"),
			expectedError: io.EOF,
		},
		{
			name:          "array header with no elements",
			input:         []byte("*1\r\n"),
			expectedError: io.EOF,
		},
		{
			name:          "nested array missing inner element",
			input:         []byte("*1\r\n*2\r\n:1\r\n"),
			expectedError: io.EOF,
		},
		{
			name:          "array length of -2",
			input:         []byte("*-2\r\n"),
			expectedError: ErrInvalidLength,
		},
		{
			name:          "array large negative length",
			input:         []byte("*-100\r\n+OK\r\n"),
			expectedError: ErrInvalidLength,
		},
		{
			name:          "nested array with invalid length",
			input:         []byte("*1\r\n*-2\r\n"),
			expectedError: ErrInvalidLength,
		},
		{
			name:          "array element with invalid bulk string length",
			input:         []byte("*1\r\n$-2\r\n"),
			expectedError: ErrInvalidLength,
		},
		{
			name:          "array length exceeds limit",
			input:         []byte("*9223372036854775807\r\n"),
			expectedError: ErrTooLarge,
		},
		{
			name:          "non-numeric array length",
			input:         []byte("*x\r\n"),
			expectedError: strconv.ErrSyntax,
		},
	}

	// Loop through each test case
	for _, tt := range tests {

		// Create subtest with given name
		t.Run(tt.name, func(t *testing.T) {
			// Wrap raw byte input in a *bufio.Reader
			reader := bufio.NewReader(bytes.NewReader(tt.input))

			val, err := Parse(reader)

			// Assert error
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Expected error: %v, actual error: %v", tt.expectedError, err)
			}

			// Assert parsed value
			if tt.expectedError == nil {
				if tt.expectedVal.Type != val.Type {
					t.Errorf("Expected type: %q, actual type: %q", tt.expectedVal.Type, val.Type)
				}

				if val.Str != tt.expectedVal.Str {
					t.Errorf("Expected string content: %q, actual string content: %q", tt.expectedVal.Str, val.Str)
				}

				if val.Int != tt.expectedVal.Int {
					t.Errorf("Expected integer value: %d, actual integer value: %d", tt.expectedVal.Int, val.Int)
				}

				if !reflect.DeepEqual(val.Array, tt.expectedVal.Array) {
					t.Errorf("Expected array: %+v, actual array: %+v", tt.expectedVal.Array, val.Array)
				}

				if val.IsNull != tt.expectedVal.IsNull {
					t.Errorf("Expected IsNull value: %t, actual IsNull value: %t", tt.expectedVal.IsNull, val.IsNull)
				}

				// Assert only the parsed value was consumed from the reader
				remaining, _ := io.ReadAll(reader)
				if string(remaining) != tt.expectedRemaining {
					t.Errorf("Expected remaining bytes: %q, actual remaining bytes: %q", tt.expectedRemaining, remaining)
				}
			}

		})
	}

}
