package redis

import (
	"github.com/codecrafters-io/redis-starter-go/internal/proto"
)

func (s *Server) HandlePing(cmd proto.Command, w *proto.Writer) error {
	return w.WriteBulkString("PONG")
}

func (s *Server) HandleEcho(cmd proto.Command, w *proto.Writer) error {
	ecmd := cmd.(proto.Echo)
	return w.WriteBulkString(ecmd.Msg())
}
