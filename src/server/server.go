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
	err = aof.LoadAOF(func(command string) error {
		response := HandleCommand(command)
		if response == protocol.ErrorValue("ERR") {
			return fmt.Errorf("error from server while handling command: %s", command)
		}
		return nil
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

		fullCommand := command + " " + strings.Join(args, " ")
		response := HandleCommand(fullCommand)
		conn.Write([]byte(marshalResponse(response)))

		if command != "" {
			aof.AppendCommand(fullCommand)
		}
	}
}

func parseORSPCommand(value protocol.ORSPValue) (string, []string, error) {
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
	args := make([]string, len(arrayValue)-1)
	for i, arg := range arrayValue[1:] {
		argStr, ok := arg.(protocol.BulkStringValue)
		if !ok {
			return "", nil, fmt.Errorf("invalid argument format: arguments must be bulk strings")
		}
		args[i] = string(argStr)
	}

	return command, args, nil
}

func marshalResponse(response string) string {
	// This is a simple implementation. You might want to enhance this based on your specific needs.
	if strings.HasPrefix(response, "ERROR:") {
		return protocol.ErrorValue(strings.TrimPrefix(response, "ERROR: ")).Marshal()
	}
	return protocol.SimpleStringValue(response).Marshal()
}
