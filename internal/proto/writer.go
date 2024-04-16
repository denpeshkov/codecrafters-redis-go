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

func (w *Writer) WriteBulkString(v string) error {
	fmt.Fprintf(w.w, "$%d\r\n%s\r\n", len(v), v)
	return w.w.Flush()
}
