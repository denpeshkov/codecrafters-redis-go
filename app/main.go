package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/redis"
)

func main() {
	s := redis.NewServer("0.0.0.0:6379")
	if err := s.Start(); err != nil {
		fmt.Println("Error starting server: ", err)
		os.Exit(1)
	}
}
