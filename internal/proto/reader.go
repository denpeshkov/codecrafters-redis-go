package proto

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/errs"
)

type Reader struct {
	r *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func (r *Reader) Read() (v any, err error) {
	defer errs.Wrap(&err, "Read")

	t, err := r.peekType()
	if err != nil {
		return nil, errors.New("malformed input: no type byte")
	}

	switch t {
	case String:
	case Error:
	case Integer:
	case BulkString:
		return r.ReadBulkString()
	case Array:
		return r.ReadArray()
	case Null:
	case Boolean:
	case Double:
	case BigNumber:
	case BulkError:
	case VerbString:
	case Map:
	case Set:
	case Push:
	}
	return nil, fmt.Errorf("unimplemented for type %q", t)
}

func (r *Reader) ReadBulkString() (bstr string, err error) {
	defer errs.Wrap(&err, "Reader.ReadBulkString")

	s, err := r.r.ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Printf("%q\n", s)

	// <length>
	l, err := strconv.Atoi(s[1 : len(s)-2])
	if err != nil {
		return "", err
	}

	// <data>
	b := make([]byte, l+2) // +2 for \r\n
	_, err = io.ReadFull(r.r, b)
	if err != nil {
		return "", err
	}
	b = b[:l]
	return string(b), nil
}

// *<number-of-elements>\r\n<element-1>...<element-n>
func (r *Reader) ReadArray() (arr []any, err error) {
	defer errs.Wrap(&err, "Reader.ReadArray")

	s, err := r.r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Printf("%q\n", s)

	// <number-of-elements>
	n, err := strconv.Atoi(s[1 : len(s)-2])
	if err != nil {
		return nil, err
	}
	arr = make([]any, n)

	// <element>s
	for i := range n {
		el, err := r.Read()
		if err != nil {
			return nil, err
		}
		arr[i] = el
	}

	return arr, nil
}

func (r *Reader) peekType() (rtype byte, err error) {
	b, err := r.r.Peek(1)
	return b[0], err
}
