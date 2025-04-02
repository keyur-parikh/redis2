package main

import (
	"errors"
	"fmt"
	"github.com/keyur-parikh/redis2/internal/definitions"
	"github.com/keyur-parikh/redis2/internal/function_mapper"
	"github.com/keyur-parikh/redis2/internal/parser"
	"github.com/keyur-parikh/redis2/internal/writer"
	"io"
	"net"
	"os"
)

func handleConnection(conn net.Conn, channel chan<- definitions.CommandInfo) {
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
		fmt.Println("Came out of read buffer appendations")
		// Send the readBuffer for parsing
		request, complete, err := parser.ValidCheckParsing(&readBuffer)
		fmt.Println("came out of valid check parsing")
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
		if len(request) == 0 {
			fmt.Printf("[%s] Parsed command was empty, skipping.\n", connID)
			continue
		}
		commandInfo := definitions.CommandInfo{
			ParsedCommand: request,
			Connection:    conn,
		}
		fmt.Println("Sending information to the database worker")
		channel <- commandInfo

	}
}

func handleDatabase(channel <-chan definitions.CommandInfo) {
	fmt.Println("Made it to handle database")
	stringToKeys := make(map[string]definitions.RedisKey)
	KVStore := make(map[definitions.RedisKey]definitions.RedisValue)
	for {
		commandInfo := <-channel
		command := commandInfo.ParsedCommand
		connection := commandInfo.Connection
		// Now I want to know what function to call
		function, err := function_mapper.FunctionMapper(command)
		if err != nil {
			// Eventually we will send the error via the connection
			fmt.Println(err)
		}
		response, err := function(command[1:], KVStore, stringToKeys)
		if err != nil {
			// Again, we can later pass on this
			// For now we send the incorrect command
			_, err := connection.Write([]byte("$-1\r\n"))
			if err != nil {
				fmt.Println("couldn't write back")
			}
		}
		if err == nil {
			if response == nil {
				//Write the Success Command and send it over
				_, err := connection.Write([]byte("+OK\\n"))
				if err != nil {
					fmt.Println("couldn't write back")
				}
			} else {
				byteResponse := writer.ArrayResponseWriter(response)
				_, err := connection.Write(byteResponse)
				if err != nil {
					fmt.Println("couldn't write back")
				}
			}

		}

	}

	// Now I want the slice of strings that I will write

}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	// here is where I should create the map
	// Uncomment this block to pass the first stage
	//
	channel := make(chan definitions.CommandInfo)
	defer close(channel)
	go handleDatabase(channel)
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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Accepted connection with ", conn.RemoteAddr().String())
		go handleConnection(conn, channel)
	}
}
