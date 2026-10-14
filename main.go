package main

import (
	"fmt"
	"github.com/ruinabeshima/mini-redis/handler"
	"github.com/ruinabeshima/mini-redis/store"
	"log"
	"net"
	"os"
)

func main() {

	// Redis server runs on PORT 6379
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Println("Connection failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Server started on PORT 6379")
	defer ln.Close()

	// Define handler and KV Store
	s := store.New()
	h := handler.New(s)

	for {

		// Client connection
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Server failed to establish connection: ", err)
			continue
		}

		go handleConnection(conn, h)
	}

}
