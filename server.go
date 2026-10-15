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

func handleConnection(conn net.Conn, h *handler.Handler) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String() // Remote network address
	log.Printf("Client connected: %s\n", remoteAddr)

	reader := bufio.NewReader(conn)
	writer := resp.NewWriter(conn)
	defer writer.Flush()

	for {
		// Parse byte stream
		val, err := resp.Parse(reader)

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("Client connection closed gracefully: %s\n", remoteAddr)
				return
			}

			log.Printf("Parse error: %v\n", err) // Log to server

			if err := writer.WriteValue(resp.Value{Type: resp.SimpleError, Str: "ERR " + err.Error()}); err != nil {
				log.Printf("Write error to %s: %v\n", remoteAddr, err)
			}
			return
		}

		// Handle command
		comm, err := handler.ParseCommand(val)
		var reply resp.Value
		if err != nil {
			log.Printf("Handler error: %v\n", err)
			reply = resp.Value{Type: resp.SimpleError, Str: "ERR " + err.Error()}
		} else {
			reply = h.ExecuteCommand(comm)
		}

		if err := writer.WriteValue(reply); err != nil {
			log.Printf("Write error to %s: %v\n", remoteAddr, err)
			return
		}

		// Flush if nothing else is pipelined
		if reader.Buffered() == 0 {
			if err := writer.Flush(); err != nil {
				log.Printf("Write error to %s: %v\n", remoteAddr, err)
				return
			}
		}
	}
}
