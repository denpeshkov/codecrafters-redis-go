package proto

import (
	"fmt"
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/errs"
)

type Command struct {
	Name string
	Args []string
}

type CommandReader struct {
	r *Reader
}

func NewCommandReader(r io.Reader) *CommandReader {
	return &CommandReader{r: NewReader(r)}
}

// *<number-of-elements>\r\n<bulk string>...<bulk string>
func (r *CommandReader) Read() (cmd *Command, err error) {
	defer errs.Wrap(&err, "CommandReader.Read")

	s, err := r.r.r.ReadString('\n')
	if err != nil {
		return cmd, err
	}
	fmt.Printf("%q\n", s)

	// <number-of-elements>
	n, err := strconv.Atoi(s[1 : len(s)-2])
	if err != nil {
		return cmd, err
	}
	arr := make([]string, n)

	// <bulk string>s
	for i := range n {
		str, err := r.r.ReadBulkString()
		if err != nil {
			return cmd, err
		}
		fmt.Printf("%q\n", str)
		arr[i] = str
	}

	return &Command{Name: arr[0], Args: arr[1:]}, nil
}
