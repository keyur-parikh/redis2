package main

import (
	"errors"
	"fmt"
	"github.com/keyur/redis2/internal/definitions"
	"github.com/keyur/redis2/internal/parser"
	"io"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	requestContext := &definitions.RequestContext{Connection: conn, KVStore: make(map[string]string)}
	connID := requestContext.Connection.RemoteAddr().String() // Unique identifier for this connection
	fmt.Printf("[%s] New connection\n", connID)
	defer func() {
		conn.Close()
		fmt.Printf("[%s] Connection closed\n", connID)
	}()

	readBuffer := make([]byte, 0)
	for {
		fmt.Printf("[%s] Reading...\n", connID)
		tempBuffer := make([]byte, 128)
		n, err := requestContext.Connection.Read(tempBuffer)

		if err != nil {
			if errors.Is(err, io.EOF) {
				// logic when client closes connection
				fmt.Println("Connection closed by client")
				return
			}
			fmt.Printf("[%s] Failed to read: %v\n", connID, err)
			return // This should exit the function
		}

		fmt.Printf("Read %d bytes: %q\n", n, tempBuffer[:n])
		readBuffer = append(readBuffer, tempBuffer[:n]...)
		// Send the readBuffer for parsing
		complete, err := parser.ValidCheckParsing(&readBuffer, requestContext)
		if err != nil {
			if complete == 2 {
				fmt.Println(err)
				return
			}
			if complete == 1 {
				fmt.Println(err)
				continue
			}

		}
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	// here is where I should create the map
	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	fmt.Println("Binded to Port ")
	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(l)

	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Started Listening to Conections")
	for {
		fmt.Println("Started Debugging")
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Accepted connection with ", conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}
