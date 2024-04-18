package proto

import (
	"bufio"
	"fmt"
	"io"
)

type Writer struct {
	w *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

func (w *Writer) Flush() error {
	return w.w.Flush()
}

func (w *Writer) WriteBulkString(v string) (int, error) {
	return fmt.Fprintf(w.w, "$%d\r\n%s\r\n", len(v), v)
}

func (w *Writer) WriteNil() (int, error) {
	return w.w.WriteString("$-1\r\n")
}
