package main

import (
	"github.com/Digital-Voting-Team/node-connector/pkg/httpserver"
	"log"
)

func main() {
	server := httpserver.NewServer()

	err := server.Echo.Start(":8080")
	if err != nil {
		log.Panicln(err)
	}
}
