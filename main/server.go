/*
	Goroutine to handle client TCP connections
*/

package main

import (
	"errors"
	"io"
	"log"
	"mini-redis/resp"
	"net"
	"strings" // Temporary package for output command string
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String() // Remote network address
	log.Printf("Client connected: %s\n", remoteAddr)

	// Temporary, immediate buffer to store incoming byte stream
	readBuffer := make([]byte, 1024)

	// Persistent buffer to store FULL byte stream as bytes come in separate chunks
	var streamBuffer []byte

	for {

		// Read number of bytes from immediate buffer
		numBytes, err := conn.Read(readBuffer)
		if err != nil {
			if errors.Is(err, io.EOF) { // Ctrl + C handler
				log.Printf("Client connection closed gracefully: %s\n", remoteAddr)
			} else {
				log.Printf("Message not received: %v\n", err)
			}
			return
		}

		log.Printf("Raw Bytes Received: %q (Hex: %x)\n", readBuffer[:numBytes], readBuffer[:numBytes])

		// Append new bytes onto persistent buffer
		streamBuffer = append(streamBuffer, readBuffer[:numBytes]...)

		for len(streamBuffer) > 0 {
			// Parse byte stream
			val, bytesConsumed, err := resp.Parse(streamBuffer)

			if err != nil {
				if errors.Is(err, resp.ErrIncomplete) {
					log.Println("Incomplete network read")
					break
				} else {
					log.Printf("Parse error: %v\n", err)
					return
				}
			}

			// Slice off parsed bytes to advance streamBuffer
			streamBuffer = streamBuffer[bytesConsumed:]

			// Handle commands in arrays
			comm := ""
			if val.Type == '*' && len(val.Array) > 0 {
				for i := 0; i < len(val.Array); i++ {
					comm += val.Array[i].Str + " "
				}
			}
			log.Printf("Command: %s\n", strings.TrimSpace(comm))
		}

		// Same reply regardless of input
		conn.Write([]byte("+PONG\r\n"))
	}
}
