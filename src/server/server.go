package server

import (
	"bufio"
	"fmt"
	"net"
	"orion/src/aof"
	"strings"
)

// StartServer initializes the TCP server
func StartServer(port string) {
	// Initialize AOF
	err := aof.InitAOF()
	if err != nil {
		fmt.Println("Error initializing AOF:", err)
		return
	}

	// Load AOF to restore state
	fmt.Println("Loading AOF data...")
	err = aof.LoadAOF(func(command string) error {
		response := HandleCommand(command)
		if response == "ERR" {
			// Handle error case based on the response
			return fmt.Errorf("error from server while handling command: %s", command)
		}
		return nil // Or return the error from server.HandleCommand
	})
	if err != nil {
		fmt.Println("Error loading AOF file:", err)
		return
	}
	fmt.Println("AOF data loaded successfully.")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is running and listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		input = strings.TrimSpace(input)
		response := HandleCommand(input)
		conn.Write([]byte(response + "\n"))
	}
}
