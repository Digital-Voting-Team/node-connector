package main

import (
	"github.com/Digital-Voting-Team/node-connector/pkg/httpserver"
	"log"
	"os"
)

func main() {
	server := httpserver.NewServer()
	port := os.Getenv("PORT")

	err := server.Echo.Start(":" + port)
	if err != nil {
		log.Panicln(err)
	}
}
