/*
	Functions to parse each individual data type in RESP2
	Each function returns the parsed data value, number of bytes processed, and error message if applicable
*/

package resp

import (
	"fmt"
	"strconv"
)

func parseSimpleString(data []byte) (string, int, error) {

	// Retrieve command slice
	slice, err := readLine(data, 1)
	if err != nil {
		return "", 0, fmt.Errorf("%w\n", err)
	}

	// Number of bytes processed
	consumed := 1 + len(slice) + 2

	return string(slice), consumed, nil
}

func parseSimpleError(data []byte) (string, int, error) {

	// Retrieve command slice
	slice, err := readLine(data, 1)
	if err != nil {
		return "", 0, fmt.Errorf("%w\n", err)
	}

	consumed := 1 + len(slice) + 2
	return string(slice), consumed, nil
}

// int: Integer, int: number of bytes consumed
func parseInteger(data []byte) (int, int, error) {

	// Retrieve command slice
	slice, err := readLine(data, 1)
	if err != nil {
		return 0, 0, fmt.Errorf("%w\n", err)
	}

	// Convert bytes to string, then parse to int
	num, err := strconv.Atoi(string(slice))
	if err != nil {
		return 0, 0, fmt.Errorf("%w\n", err)
	}

	consumed := 1 + len(slice) + 2
	return num, consumed, nil
}

// 　Bool return value is for isNull (null bulk string)
func parseBulkString(data []byte) (string, bool, int, error) {

	// Retrieve string length and convert to int
	length, err := readLine(data, 1)
	if err != nil {
		return "", false, 0, fmt.Errorf("%w\n", err)
	}
	intLength, err := strconv.Atoi(string(length))
	if err != nil {
		return "", false, 0, fmt.Errorf("%w\n", err)
	}

	// Negative lengths 
	if intLength < -1 {
		return "", false, 0, ErrInvalidLength
	}

	// Null bulk strings (-1)
	if intLength == -1 {
		return "", true, 5, nil
	}

	// Calculate where the payload starts and ends
	bulkStart := 1 + len(length) + 2
	bulkEnd := bulkStart + intLength
	if bulkEnd+2 > len(data) {
		return "", false, 0, ErrIncomplete
	}

	// Slice the payload directly, and verify CRLF after
	bulkBytes := data[bulkStart:bulkEnd]
	if !is_CRLF(data, bulkEnd) {
		return "", false, 0, ErrIncomplete
	}

	consumed := 1 + len(length) + 2 + intLength + 2
	return string(bulkBytes), false, consumed, nil
}

func parseArray(data []byte) (Value, int, error) {

	// Get length of array
	length, err := readLine(data, 1)
	if err != nil {
		return Value{}, 0, fmt.Errorf("%w\n", err)
	}
	intLength, err := strconv.Atoi(string(length))
	if err != nil {
		return Value{}, 0, fmt.Errorf("%w\n", err)
	}

	// Invalid length 
	if intLength < -1 {
		return Value{}, 0, ErrInvalidLength
	}

	// Null array
	if intLength == -1 {
		return Value{Type: array, IsNull: true}, 5, nil
	}

	// Calculate initial offset
	offset := 1 + len(length) + 2
	elements := make([]Value, intLength)

	//　Recursively parse each child element
	for i := 0; i < intLength; i++ {
		// Ran out of data before all elements arrived 
		if offset >= len(data) {
			return Value{}, 0, ErrIncomplete
		}

		val, consumed, err := Parse(data[offset:])
		if err != nil {
			return Value{}, 0, err
		}

		elements[i] = val
		offset += consumed
	}

	return Value{Type: array, Array: elements}, offset, nil
}
