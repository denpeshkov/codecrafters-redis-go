package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/errs"
	"github.com/codecrafters-io/redis-starter-go/internal/proto"
)

func (s *Server) HandlePing(cmd *proto.Command, w *proto.Writer) error {
	w.WriteBulkString("PONG")
	return w.Flush()
}

func (s *Server) HandleEcho(cmd *proto.Command, w *proto.Writer) error {
	w.WriteBulkString(cmd.Args[0])
	return w.Flush()
}

func (s *Server) HandleSet(cmd *proto.Command, w *proto.Writer) (err error) {
	defer errs.Wrap(&err, "HandleSet")
	type set struct {
		key string
		val string
		exp time.Time
	}

	setCmd := set{key: cmd.Args[0], val: cmd.Args[1]}
	opts := cmd.Args[2:]
	for len(opts) > 0 {
		switch opt := opts[0]; opt {
		case "px":
			exp, err := strconv.Atoi(opts[1])
			if err != nil {
				return fmt.Errorf("malformed PX value: %w", err)
			}
			setCmd.exp = time.Now().Add(time.Duration(exp) * time.Millisecond)
			opts = opts[2:]
		default:
			return fmt.Errorf("unsupported SET option: %v", opt)
		}
	}

	s.mu.Lock()
	s.db[setCmd.key] = setCmd.val
	delete(s.ttl, setCmd.key)
	if !setCmd.exp.IsZero() {
		s.ttl[setCmd.key] = setCmd.exp
	}
	s.mu.Unlock()

	w.WriteBulkString("OK")
	return w.Flush()
}

func (s *Server) HandleGet(cmd *proto.Command, w *proto.Writer) error {
	key := cmd.Args[0]

	s.mu.Lock()
	if exp, ok := s.ttl[key]; ok && time.Now().After(exp) {
		s.deleteKey(key)
	}
	v, ok := s.db[key]
	s.mu.Unlock()

	if !ok {
		w.WriteNil()
	} else {
		w.WriteBulkString(v)
	}
	return w.Flush()
}
