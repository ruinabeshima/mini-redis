package main

import (
	"fmt"
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
		bytes, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Message not received")
			continue 
		}

		// Convert to message and check for connection exit 
		message := string(buffer[:bytes])
		if message == "exit\n" || message == "exit\r\n" {
			fmt.Printf("Connection closed: %s\n", remoteAddr)
			return 
		}

		// Same reply regardless of input 
		conn.Write([]byte("+PONG\r\n"))
	}
}
