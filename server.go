/*
	Goroutine to handle client TCP connections
*/

package main

import (
	"bufio"
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

	reader := bufio.NewReader(conn)

	for {
		// Parse byte stream
		val, err := resp.Parse(reader)

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Client connection closed gracefully: %s\n", remoteAddr)
				return
			}

			log.Printf("Parse error: %v\n", err)

			// Send error to client
			if err := writeError(conn, err.Error()); err != nil {
				log.Printf("Write error to %s: %v\n", remoteAddr, err)
			}

			return
		}

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
