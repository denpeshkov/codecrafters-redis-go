package proto

import (
	"fmt"
	"io"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/errs"
)

type Command interface {
	Name() string
}

type commandRequest []string

func (c commandRequest) Name() string {
	return c[0]
}
func (c commandRequest) Args() []string {
	return c[1:]
}

type CommandReader struct {
	r *Reader
}

func NewCommandReader(r io.Reader) *CommandReader {
	return &CommandReader{r: NewReader(r)}
}

func (r *CommandReader) Read() (cmd Command, err error) {
	defer errs.Wrap(&err, "CommandReader.Read")

	cmdr, err := r.read()
	if err != nil {
		return nil, err
	}

	switch cmdr.Name() {
	case "ping":
		return NewPing(), nil
	case "echo":
		return NewEcho(cmdr.Args()[0]), nil
	}
	return nil, fmt.Errorf("unimplemented command: %q", cmdr.Name())
}

// *<number-of-elements>\r\n<bulk string>...<bulk string>
func (r *CommandReader) read() (cmd commandRequest, err error) {
	defer errs.Wrap(&err, "CommandReader.read")

	s, err := r.r.r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	fmt.Printf("%q\n", s)

	// <number-of-elements>
	n, err := strconv.Atoi(s[1 : len(s)-2])
	if err != nil {
		return nil, err
	}
	arr := make([]string, n)

	// <bulk string>s
	for i := range n {
		str, err := r.r.ReadBulkString()
		if err != nil {
			return nil, err
		}
		fmt.Printf("%q\n", str)
		arr[i] = str
	}

	return arr, nil
}

type Ping struct{}

func NewPing() Ping {
	return Ping{}
}
func (p Ping) Name() string {
	return "ping"
}

type Echo struct {
	msg string
}

func NewEcho(msg string) Echo {
	return Echo{msg: msg}
}
func (e Echo) Name() string {
	return "echo"
}
func (e Echo) Msg() string {
	return e.msg
}
