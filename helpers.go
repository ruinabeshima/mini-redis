package main

import "net"

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
