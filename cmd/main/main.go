package main

import (
	"log"
	"node-connector/httpserver"
)

func main() {
	server := httpserver.NewServer()

	err := server.Echo.Start(":8080")
	if err != nil {
		log.Panicln(err)
	}
}
