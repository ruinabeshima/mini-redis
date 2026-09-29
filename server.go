/*
	Goroutine to handle client TCP connections
*/

package main

import (
	"errors"
	"github.com/ruinabeshima/mini-redis/handler"
	"github.com/ruinabeshima/mini-redis/resp"
	"io"
	"log"
	"net"
)

// Helper to write raw byte slices to net.Conn
func writeBytes(conn net.Conn, data []byte) error {
	_, err := conn.Write(data)
	return err
}

// Helper to write RESP formatted error strings
func writeError(conn net.Conn, errMsg string) error {
	_, err := conn.Write([]byte("-ERR " + errMsg + "\r\n"))
	return err
}

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

					// Send error to client
					if err := writeError(conn, err.Error()); err != nil {
						log.Printf("Write error to %s: %v\n", remoteAddr, err)
					}

					return
				}
			}

			// Complete network read: Slice off parsed bytes to advance streamBuffer
			streamBuffer = streamBuffer[bytesConsumed:]

			// Handle command
			comm, err := handler.ParseCommand(val)
			if err != nil {
				log.Printf("Handler error: %v\n", err)

				// Send error to client
				if err := writeError(conn, err.Error()); err != nil {
					log.Printf("Write error to %s: %v\n", remoteAddr, err)
					return
				}
			} else {
				returnBytes := handler.ExecuteCommand(comm)

				// Send bytes to client
				if err := writeBytes(conn, returnBytes); err != nil {
					log.Printf("Write error to %s: %v\n", remoteAddr, err)
					return
				}
			}
		}
	}
}
