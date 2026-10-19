/*
	fmt.Fprintf formats directly on the bufio.Writer,
	meaning there is no intermediate object involved.
	This is unlike fmt.Appendf() which builds a complete byte
	slice on the heap, and []byte(fmt.Sprintf()) which builds a
	string and then copies it in a fresh byte slice.
*/

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
	_, err := fmt.Fprintf(buf.w, "+%s\r\n", v.Str)
	return err
}

func (buf *Writer) writeSimpleError(v Value) error {
	_, err := fmt.Fprintf(buf.w, "-%s\r\n", v.Str)
	return err
}

func (buf *Writer) writeInteger(v Value) error {
	_, err := fmt.Fprintf(buf.w, ":%d\r\n", v.Int)
	return err
}

func (buf *Writer) writeBulkString(v Value) error {
	// Null bulk string
	if v.IsNull {
		_, err := buf.w.WriteString("$-1\r\n")
		return err
	}

	_, err := fmt.Fprintf(buf.w, "$%d\r\n%s\r\n", len(v.Str), v.Str)
	return err
}

func (buf *Writer) writeArray(v Value) error {
	arr := v.Array
	if v.IsNull {
		_, err := buf.w.WriteString("*-1\r\n") // Null array
		return err
	}
	if _, err := fmt.Fprintf(buf.w, "*%d\r\n", len(arr)); err != nil {
		return err
	}

	for i := 0; i < len(arr); i++ {
		if err := buf.WriteValue(arr[i]); err != nil {
			return err
		}
	}
	return nil
}
