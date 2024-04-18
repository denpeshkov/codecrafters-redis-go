package redis

import (
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

func (s *Server) HandleSet(cmd *proto.Command, w *proto.Writer) error {
	type Set struct {
		Key   string
		Value string
	}
	setCmd := Set{Key: cmd.Args[0], Value: cmd.Args[1]}
	s.mu.Lock()
	s.db[setCmd.Key] = setCmd.Value
	s.mu.Unlock()

	w.WriteBulkString("OK")
	return w.Flush()
}

func (s *Server) HandleGet(cmd *proto.Command, w *proto.Writer) error {
	key := cmd.Args[0]

	s.mu.RLock()
	v, ok := s.db[key]
	s.mu.RUnlock()
	if !ok {
		w.WriteNil()
	} else {
		w.WriteBulkString(v)
	}
	return w.Flush()
}
