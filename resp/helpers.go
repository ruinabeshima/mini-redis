/*
	Helper functions to verify CRLF \r\n and read lines of data
*/

package resp

import "errors"

func is_CRLF(byteArray []byte, pointer int) bool {

	// Pointer out of bounds
	if pointer+1 >= len(byteArray) || pointer < 0 {
		return false
	}

	// Check if there is a Carriage Return Line Feed (\r\n)
	if byteArray[pointer] == '\r' && byteArray[pointer+1] == '\n' {
		return true
	}

	return false
}

/*
Finds \r\n and returns the slice of bytes from "start" up to (but not including) \r\n
For simple strings, simple errors, integers
*/
func readLine(data []byte, start int) ([]byte, error) {
	end := start
	for end < len(data) && !is_CRLF(data, end) {
		end += 1
	}

	// \r\n not included
	if end == len(data) {
		return nil, errors.New("CRLF not included")
	}

	return data[start:end], nil
}
