package resp

import (
	"bufio"
	"fmt"
	"io"
)

type Writer struct {
	w *bufio.Writer
}

func NewWriter(r io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(r)}
}

func (buf *Writer) WriteValue(v Value) error {
	// Obtain the data type
	typeByte := v.Type

	switch typeByte {
	case SimpleString:
		return buf.writeSimpleString(v)
	case SimpleError:
		return buf.writeSimpleError(v)
	case Integer:
		return buf.writeInteger(v)
	case BulkString:
		return buf.writeBulkString(v)
	case TypeArray:
		return buf.writeArray(v)
	default:
		return ErrUnknownType
	}
}

/*
With bufio, data accumulates in RAM and is not output
to avoid expensive system calls.
Flush() forces the writer to empty its contents.
*/
func (buf *Writer) Flush() error {
	return buf.w.Flush()
}

func (buf *Writer) writeSimpleString(v Value) error {
	s := v.Str
	data := []byte(fmt.Sprintf("+%s\r\n", s))

	// Writes a byte slice to the buffer
	_, err := buf.w.Write(data)
	return err
}

func (buf *Writer) writeSimpleError(v Value) error {
	s := v.Str
	data := []byte(fmt.Sprintf("-%s\r\n", s))

	_, err := buf.w.Write(data)
	return err
}

func (buf *Writer) writeInteger(v Value) error {
	num := v.Int
	data := []byte(fmt.Sprintf(":%d\r\n", num))

	_, err := buf.w.Write(data)
	return err
}

func (buf *Writer) writeBulkString(v Value) error {
	s := v.Str
	var data []byte
	if v.IsNull {
		data = []byte("$-1\r\n") // Null bulk string
	} else {
		data = []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
	}

	_, err := buf.w.Write(data)
	return err
}

func (buf *Writer) writeArray(v Value) error {
	arr := v.Array
	if v.IsNull {
		_, err := buf.w.WriteString("*-1\r\n") // Null array 
		return err
	}
	if _, err := buf.w.WriteString(fmt.Sprintf("*%d\r\n", len(arr))); err != nil {
		return err
	}

	for i := 0; i < len(arr); i++ {
		if err := buf.WriteValue(arr[i]); err != nil {
			return err
		}
	}
	return nil
}
