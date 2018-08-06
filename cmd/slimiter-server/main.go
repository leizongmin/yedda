package main

import (
	"github.com/leizongmin/simple-limiter-service/server"
	"log"
)

func main() {
	s, err := server.NewServer(server.Options{
		EnableLog: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	s.Loop()
}
