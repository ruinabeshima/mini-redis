/*
	Goroutine to handle client TCP connections
*/

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

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
		fmt.Println("Message: ", byteMessage)

		// Same reply regardless of input
		conn.Write([]byte("+PONG\r\n"))
	}
}
