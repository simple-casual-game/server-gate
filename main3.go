//go:build test
// +build test

package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {

	ipaddr := "127.0.0.1" + ":" + "8001"

	c, err := net.Listen("tcp", ipaddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// session

	conn, err := c.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from: " + remoteAddr)

	connection := &Connection{}
	connection.BufferReader = bufio.NewReaderSize(conn, 1024)
	connection.BufferWriter = bufio.NewWriterSize(conn, 1024)

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	for {

		length, err := connection.BufferReader.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Disconned from ", remoteAddr)
				break
			} else {
				fmt.Println("Error reading:", err.Error())
				break
			}
		}
		connection.OnData(buf[:length], length)

	}
	// Close the connection when you're done with it.
	conn.Close()

}
