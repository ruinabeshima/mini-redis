package resp

import (
	"bufio"
	"errors"
)

/*
Helper function to read line of data until trailing \n, and verify that it ends with \r\n
*/
func readLine(r *bufio.Reader) ([]byte, error) {
	// Read up to trailing \n
	line, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// Verify that the line ends in \r\n
	if len(line) < 2 || line[len(line)-2] != '\r' {
		return nil, ErrNoCRLF
	}

	return line[:len(line)-2], nil
}
