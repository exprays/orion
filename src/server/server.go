package server

import (
	"bufio"
	"fmt"
	"net"
	"orion/src/aof"
	"orion/src/protocol"
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
	err = aof.LoadAOF(func(command protocol.ArrayValue) error {
		response := HandleCommand(command)
		if errValue, ok := response.(protocol.ErrorValue); ok {
			fmt.Printf("Warning: Error handling command %v: %s\n", command, string(errValue))
			// Continue loading instead of returning an error
			return nil
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error loading AOF: %v\n", err)
		// Consider whether you want to continue starting the server or exit here
	} else {
		fmt.Println("AOF data loaded successfully.")
	}

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
		value, err := protocol.Unmarshal(reader)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		command, args, err := parseORSPCommand(value)
		if err != nil {
			response := protocol.ErrorValue(err.Error())
			conn.Write([]byte(response.Marshal()))
			continue
		}

		fullArgs := append([]protocol.ORSPValue{protocol.BulkStringValue(command)}, args...)
		response := HandleCommand(fullArgs)
		conn.Write([]byte(response.Marshal()))

		if command != "" {
			fullCommand := protocol.ArrayValue{protocol.BulkStringValue(command)}
			for _, arg := range args {
				fullCommand = append(fullCommand, arg)
			}
			aof.AppendCommand(fullCommand) // Convert and pass as ArrayValue
		}
	}
}

func parseORSPCommand(value protocol.ORSPValue) (string, []protocol.ORSPValue, error) {
	arrayValue, ok := value.(protocol.ArrayValue)
	if !ok {
		return "", nil, fmt.Errorf("invalid command format: expected array")
	}

	if len(arrayValue) == 0 {
		return "", nil, fmt.Errorf("empty command")
	}

	commandStr, ok := arrayValue[0].(protocol.BulkStringValue)
	if !ok {
		return "", nil, fmt.Errorf("invalid command format: command must be a bulk string")
	}

	command := strings.ToUpper(string(commandStr))
	args := arrayValue[1:]

	return command, args, nil
}

// func parseStringCommand(command string) []protocol.ORSPValue {
// 	parts := strings.Fields(command)
// 	args := make([]protocol.ORSPValue, len(parts))
// 	for i, part := range parts {
// 		args[i] = protocol.BulkStringValue(part)
// 	}
// 	return args
// }
