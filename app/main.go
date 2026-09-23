package main

import (
	"bytes"
	"errors"
	"io"
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
	defer ln.Close()

	for {

		// Client connection
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Server failed to establish connection: ", err)
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Remote network address
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s\n", remoteAddr)

	// Buffer to store data
	buffer := make([]byte, 1024)

	for {

		// Read number of bytes
		numBytes, err := conn.Read(buffer)
		if err != nil {
			// Handler for standard nc client pressing Ctrl + C (Terminal intercepts lcoally and kills client process immediately)
			if errors.Is(err, io.EOF) {
				log.Printf("Client connection closed gracefully: %s\n", remoteAddr)
			} else {
				log.Printf("Message not received: %v\n", err)
			}
			return
		}

		byteMessage := buffer[:numBytes]

		// Same reply regardless of input
		conn.Write([]byte("+PONG\r\n"))
	}
}
