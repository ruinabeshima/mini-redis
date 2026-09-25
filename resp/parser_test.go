/*
	Table driven tests for Parse() function
*/

package resp

import (
	"errors"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {

	// Declare and initialise test cases
	tests := []struct {
		name          string
		input         []byte
		expectedVal   Value
		expectedBytes int
		expectedError error
	}{
		// Simple strings
		{
			name:          "simple string",
			input:         []byte("+OK\r\n"),
			expectedVal:   Value{Type: simpleString, Str: "OK"},
			expectedBytes: 5,
		},
		{
			name:          "empty simple string",
			input:         []byte("+\r\n"),
			expectedVal:   Value{Type: simpleString, Str: ""},
			expectedBytes: 3,
		},
		{
			name:          "simple string with trailing data only consumes first value",
			input:         []byte("+OK\r\n+PONG\r\n"),
			expectedVal:   Value{Type: simpleString, Str: "OK"},
			expectedBytes: 5,
		},
		{
			name:          "simple string missing CRLF",
			input:         []byte("+OK"),
			expectedError: ErrIncomplete,
		},
		{
			name:          "simple string missing LF",
			input:         []byte("+OK\r"),
			expectedError: ErrIncomplete,
		},

		// Simple errors
		{
			name:          "simple error",
			input:         []byte("-ERR unknown command\r\n"),
			expectedVal:   Value{Type: simpleError, Str: "ERR unknown command"},
			expectedBytes: 22,
		},
		{
			name:          "simple error missing CRLF",
			input:         []byte("-ERR"),
			expectedError: ErrIncomplete,
		},

		// Integers
		{
			name:          "positive integer",
			input:         []byte(":1000\r\n"),
			expectedVal:   Value{Type: integer, Int: 1000},
			expectedBytes: 7,
		},
		{
			name:          "negative integer",
			input:         []byte(":-42\r\n"),
			expectedVal:   Value{Type: integer, Int: -42},
			expectedBytes: 6,
		},
		{
			name:          "explicit plus sign integer",
			input:         []byte(":+5\r\n"),
			expectedVal:   Value{Type: integer, Int: 5},
			expectedBytes: 5,
		},
		{
			name:          "zero integer",
			input:         []byte(":0\r\n"),
			expectedVal:   Value{Type: integer, Int: 0},
			expectedBytes: 4,
		},
		{
			name:          "integer missing CRLF",
			input:         []byte(":12"),
			expectedError: ErrIncomplete,
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
			name:          "bulk string",
			input:         []byte("$5\r\nhello\r\n"),
			expectedVal:   Value{Type: bulkString, Str: "hello"},
			expectedBytes: 11,
		},
		{
			name:          "empty bulk string",
			input:         []byte("$0\r\n\r\n"),
			expectedVal:   Value{Type: bulkString, Str: ""},
			expectedBytes: 6,
		},
		{
			name:          "bulk string containing CRLF is binary safe",
			input:         []byte("$7\r\nfoo\r\nba\r\n"),
			expectedVal:   Value{Type: bulkString, Str: "foo\r\nba"},
			expectedBytes: 13,
		},
		{
			name:          "bulk string with multi-digit length",
			input:         []byte("$11\r\nhello world\r\n"),
			expectedVal:   Value{Type: bulkString, Str: "hello world"},
			expectedBytes: 18,
		},
		{
			name:          "null bulk string",
			input:         []byte("$-1\r\n"),
			expectedVal:   Value{Type: bulkString, IsNull: true},
			expectedBytes: 5,
		},
		{
			name:          "bulk string missing length CRLF",
			input:         []byte("$5"),
			expectedError: ErrIncomplete,
		},
		{
			name:          "bulk string payload shorter than length",
			input:         []byte("$5\r\nhel"),
			expectedError: ErrIncomplete,
		},
		{
			name:          "bulk string missing trailing CRLF",
			input:         []byte("$5\r\nhello"),
			expectedError: ErrIncomplete,
		},
		{
			name:          "bulk string payload longer than length",
			input:         []byte("$3\r\nhello\r\n"),
			expectedError: ErrIncomplete,
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
			name:          "non-numeric bulk string length",
			input:         []byte("$abc\r\nhello\r\n"),
			expectedError: strconv.ErrSyntax,
		},

		// Arrays
		{
			name:          "empty array",
			input:         []byte("*0\r\n"),
			expectedVal:   Value{Type: array, Array: []Value{}},
			expectedBytes: 4,
		},
		{
			name:          "null array",
			input:         []byte("*-1\r\n"),
			expectedVal:   Value{Type: array, IsNull: true},
			expectedBytes: 5,
		},
		{
			name:  "array of bulk strings",
			input: []byte("*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n"),
			expectedVal: Value{Type: array, Array: []Value{
				{Type: bulkString, Str: "ECHO"},
				{Type: bulkString, Str: "hello"},
			}},
			expectedBytes: 25,
		},
		{
			name:  "array of mixed types",
			input: []byte("*3\r\n:1\r\n+OK\r\n$-1\r\n"),
			expectedVal: Value{Type: array, Array: []Value{
				{Type: integer, Int: 1},
				{Type: simpleString, Str: "OK"},
				{Type: bulkString, IsNull: true},
			}},
			expectedBytes: 18,
		},
		{
			name:  "nested array",
			input: []byte("*2\r\n*1\r\n:1\r\n+OK\r\n"),
			expectedVal: Value{Type: array, Array: []Value{
				{Type: array, Array: []Value{{Type: integer, Int: 1}}},
				{Type: simpleString, Str: "OK"},
			}},
			expectedBytes: 17,
		},
		{
			name:  "array with trailing data only consumes first value",
			input: []byte("*1\r\n+PING\r\n*1\r\n+PING\r\n"),
			expectedVal: Value{Type: array, Array: []Value{
				{Type: simpleString, Str: "PING"},
			}},
			expectedBytes: 11,
		},
		{
			name:          "array missing length CRLF",
			input:         []byte("*2"),
			expectedError: ErrIncomplete,
		},
		{
			name:          "array with incomplete element",
			input:         []byte("*2\r\n$4\r\nECHO\r\n$5\r\nhel"),
			expectedError: ErrIncomplete,
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
			name:          "non-numeric array length",
			input:         []byte("*x\r\n"),
			expectedError: strconv.ErrSyntax,
		},
	}

	// Loop through each test case
	for _, tt := range tests {

		// Create subtest with given name
		t.Run(tt.name, func(t *testing.T) {
			val, bytesConsumed, err := Parse(tt.input)

			// Assert error
			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Expected error: %v, actual error: %v", tt.expectedError, err)
			}

			// Assert bytes consumed
			if tt.expectedBytes != bytesConsumed {
				t.Errorf("Expected bytes consumed: %d, actual bytes consumed: %d", tt.expectedBytes, bytesConsumed)
			}

			// Assert parsed value
			if tt.expectedError == nil {
				if tt.expectedVal.Type != val.Type {
					t.Errorf("Expected type: %v, actual type: %v", tt.expectedVal.Type, val.Type)
				}

				if val.Str != tt.expectedVal.Str {
					t.Errorf("Expected string content: %q, actual string content: %q", tt.expectedVal.Str, val.Str)
				}

				if val.Int != tt.expectedVal.Int {
					t.Errorf("Expected integer value: %d, actual integer value: %d", tt.expectedVal.Int, val.Int)
				}

				if len(val.Array) != len(tt.expectedVal.Array) {
					t.Errorf("Expected array length: %d, actual array length: %d", len(tt.expectedVal.Array), len(val.Array))
				}

				if val.IsNull != tt.expectedVal.IsNull {
					t.Errorf("Expected IsNull value: %t, actual IsNull value: %t", val.IsNull, tt.expectedVal.IsNull)
				}
			}

		})
	}

}
