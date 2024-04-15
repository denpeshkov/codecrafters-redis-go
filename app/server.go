package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

func run() error {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		return fmt.Errorf("bind port 6379: %w", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("accepting connection: %w", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	b := make([]byte, 1024)
	for {
		n, err := conn.Read(b)
		if errors.Is(err, io.EOF) {
			return
		}
		if err != nil {
			fmt.Println("Reading from connection: ", err)
		}

		fmt.Printf("received: %q\n", string(b[:n]))
		conn.Write([]byte("+PONG\r\n"))
	}
}
