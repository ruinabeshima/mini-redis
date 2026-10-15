package resp

import "bufio"

type Writer struct {
	w *bufio.Writer
}

func NewWriter() *Writer {
	return &Writer
}

func (w *Writer) Write(v Value) error {
	// Obtain the data type
	typeByte := v.Type

	switch typeByte {
	case Simplestring:
		err := WriteSimpleString()
	case SimpleError:
		err := WriteSimpleError()
	case Integer:
		err := WriteInteger()
	case BulkString:
		err := WriteBulkString()
	case TypeArray:
		err := WriteArray()
	}
}

/*
With bufio, data accumulates in RAM and is not output
to avoid expensive system calls.
Flush() forces the writer to empty its contents.
*/
func (w *Writer) Flush() error {
	err := w.Flush()
	if err != nil {
		return err
	}
}

func (w *Writer) WriteSimpleString() error {

}

func (w *Writer) WriteSimpleError() error {

}

func (w *Writer) WriteInteger() error {

}

func (w *Writer) WriteBulkString() error {

}

func (w *Writer) WriteArray() error {

}
