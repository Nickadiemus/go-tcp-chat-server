package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/nickadiemus/go-tcp-chat-server/pkg/server"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
	s := server.NewServer()
	go s.Run()
	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
	defer listener.Close() // ensure connection is closed

	log.Printf("Listening on localhost:%s\n", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("could not accept connection: %s", err.Error())
			continue
		}

		// we want to handle many connections concurrently
		go s.NewClient(conn)
	}
}
