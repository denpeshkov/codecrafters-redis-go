package redis

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/internal/errs"
	"github.com/codecrafters-io/redis-starter-go/internal/proto"
)

type CommandHandler func(cmd *proto.Command, w *proto.Writer) error

type Server struct {
	addr string

	mu   sync.RWMutex
	cmds map[string]CommandHandler
	db   map[string]string
}

func NewServer(addr string) *Server {
	s := &Server{
		addr: addr,
		cmds: make(map[string]CommandHandler),
		db:   make(map[string]string),
	}

	s.Register("ping", s.HandlePing)
	s.Register("echo", s.HandleEcho)
	s.Register("set", s.HandleSet)
	s.Register("get", s.HandleGet)

	return s
}

func (s *Server) Start() (err error) {
	defer errs.Wrap(&err, "Server.Start")

	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen on address: %q: %w", s.addr, err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("accepting connection on address: %q: %w", s.addr, err)
		}
		go s.ServeConn(conn)
	}
}

func (s *Server) Register(name string, cmdh CommandHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	name = strings.ToLower(name)
	s.cmds[name] = cmdh
}

func (s *Server) ServeConn(conn net.Conn) {
	r := proto.NewCommandReader(conn)
	w := proto.NewWriter(conn)

	for {
		cmd, err := r.Read()
		if errors.Is(err, io.EOF) {
			fmt.Println("Closed connection")
			return
		}
		if err != nil {
			fmt.Println("Error reading command: ", err)
			return
		}
		s.mu.RLock()
		h, ok := s.cmds[cmd.Name]
		s.mu.RUnlock()
		if !ok {
			fmt.Printf("No handler for command: %q\n", cmd.Name)
		}
		if err := h(cmd, w); err != nil {
			fmt.Printf("Error processing command: %q: %v\n", cmd.Name, err)
		}
	}
}
