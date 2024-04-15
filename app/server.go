package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		return fmt.Errorf("bind port 6379: %w", err)
	}

	conn, err := l.Accept()
	if err != nil {
		return fmt.Errorf("accepting connection: %w", err)
	}

	b := make([]byte, 1024)
	for {
		n, err := conn.Read(b)
		if err != nil && err != io.EOF {
			return fmt.Errorf("reading from connection: %w", err)
		}
		fmt.Printf("received: %q\n", string(b[:n]))
		conn.Write([]byte("+PONG\r\n"))
	}
}
