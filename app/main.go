package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	// Redis runs on PORT 6379
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Connection failed: ", err)
		os.Exit(1)
	}
	defer ln.Close()

	// Run continuously
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection not established \n")
			continue
		}

		// Handle connection
		go handleConnection(conn)
	}

}

function handleConnection(conn net.Conn) {
	defer conn.Close() 
}
