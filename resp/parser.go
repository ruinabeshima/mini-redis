/*
	Functions to parse each individual data type in RESP2
	Each function returns the parsed data value and error message if applicable
*/

package resp

import (
	"bufio"
	"io"
	"strconv"
)

func parseSimpleString(r *bufio.Reader) (string, error) {

	// Retrieve command slice
	slice, err := readLine(r)
	if err != nil {
		return "", err
	}

	return string(slice), nil
}

func parseSimpleError(r *bufio.Reader) (string, error) {

	// Retrieve command slice
	slice, err := readLine(r)
	if err != nil {
		return "", err
	}

	return string(slice), nil
}

func parseInteger(r *bufio.Reader) (int, error) {

	// Retrieve command slice
	slice, err := readLine(r)
	if err != nil {
		return 0, err
	}

	// Convert bytes to string, then parse to int
	num, err := strconv.Atoi(string(slice))
	if err != nil {
		return 0, err
	}

	return num, nil
}

// 　Bool return value is for isNull (null bulk string)
func parseBulkString(r *bufio.Reader) (string, bool, error) {

	// Retrieve string length and convert to int
	length, err := readLine(r)
	if err != nil {
		return "", false, err
	}
	intLength, err := strconv.Atoi(string(length))
	if err != nil {
		return "", false, err
	}

	// Negative lengths
	if intLength < -1 {
		return "", false, ErrInvalidLength
	}

	// Length too large
	if intLength > maxBulkLength {
		return "", false, ErrTooLarge
	}

	// Null bulk strings (-1)
	if intLength == -1 {
		return "", true, nil
	}

	// Read exact length into a buffer of length intLength
	buf := make([]byte, intLength)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", false, err
	}

	// Read and validate trailing \r\n
	crlfBuf := make([]byte, 2)
	if _, err := io.ReadFull(r, crlfBuf); err != nil {
		return "", false, err
	}
	if crlfBuf[0] != '\r' || crlfBuf[1] != '\n' {
		return "", false, fmt.Errorf("bulk string missing CRLF ending")
	}

	return string(buf), false, nil
}

func parseArray(r *bufio.Reader) (Value, error) {

	// Get length of array
	length, err := readLine(r)
	if err != nil {
		return Value{}, err
	}
	intLength, err := strconv.Atoi(string(length))
	if err != nil {
		return Value{}, err
	}

	// Invalid length
	if intLength < -1 {
		return Value{}, ErrInvalidLength
	}

	// Null array
	if intLength == -1 {
		return Value{Type: array, IsNull: true}, nil
	}

	// Array length too large
	if intLength > maxArrayLength {
		return Value{}, ErrTooLarge
	}

	elements := make([]Value, intLength)

	//　Recursively parse each child element
	for i := 0; i < intLength; i++ {
		val, err := Parse(r)
		if err != nil {
			return Value{}, err
		}

		elements[i] = val
	}

	return Value{Type: array, Array: elements}, nil
}
