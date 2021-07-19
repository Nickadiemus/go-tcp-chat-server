package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nickadiemus/go-tcp-chat-server/pkg/server"
)

func main() {
	err := godotenv.Load()
	fmt.Println("Starting chat server...")
	s := server.NewServer()
	go s.Run()
	port := os.Getenv("PORT")
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
	defer listener.Close() // ensure connection is closed
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("could not accept connection: %s", err.Error())
			continue
		}

		// we want to handle many connections concurrent
		go s.NewClient(conn)
	}
}
