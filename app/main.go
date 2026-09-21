package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	// Redis server runs on PORT 6379
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Connection failed: ", err)
		os.Exit(1)
	}
	defer ln.Close()

	for {

		// Client connection
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection not established")
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Remote network address
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", remoteAddr)

	// Buffer to store data
	buffer := make([]byte, 1024)

	for {

		// Read number of bytes
		numBytes, err := conn.Read(buffer)
		if err != nil {
			// Handler for standard nc client pressing Ctrl + C (Terminal intercepts lcoally and kills client process immediately)
			if errors.Is(err, io.EOF) {
				fmt.Printf("Server closed gracefully: %s\n", remoteAddr)
			} else {
				fmt.Println("Message not received")
			}
			return
		}

		byteMessage := buffer[:numBytes]

		// Handle Ctrl+C byte for clients running in raw mode (input contains ASCII value 3)
		if bytes.IndexByte(byteMessage, 0x03) != -1 {

			fmt.Printf("Ctrl+C received from %s. Closing client connection.\n", remoteAddr)
			return
		}

		// Same reply regardless of input
		conn.Write([]byte("+PONG\r\n"))
	}
}
