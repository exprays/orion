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
	// Initialize logging system
	err := InitLogging()
	if err != nil {
		fmt.Println("Error initializing logging system:", err)
		return
	}
	defer CloseLogFiles()

	// Initialize AOF
	err = aof.InitAOF()
	if err != nil {
		LogError("Error initializing AOF: %v", err)
		return
	}

	// Load AOF to restore state
	LogInfo("Loading AOF data...")
	err = aof.LoadAOF(func(command protocol.ArrayValue) error {
		response := HandleCommand(command)
		if errValue, ok := response.(protocol.ErrorValue); ok {
			LogError("Warning: Error handling command %v: %s", command, string(errValue))
			// Continue loading instead of returning an error
			return nil
		}
		return nil
	})

	if err != nil {
		LogError("Error loading AOF: %v", err)
		// Consider whether you want to continue starting the server or exit here
	} else {
		LogInfo("AOF data loaded successfully.")
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		LogError("Error starting server: %v", err)
		return
	}
	defer listener.Close()

	LogInfo("Server is running and listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			LogError("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	LogInfo("New connection from %s", clientAddr)

	reader := bufio.NewReader(conn)
	for {
		value, err := protocol.Unmarshal(reader)
		if err != nil {
			LogError("Error reading input from %s: %v", clientAddr, err)
			return
		}

		command, args, err := parseORSPCommand(value)
		if err != nil {
			LogError("Invalid command from %s: %v", clientAddr, err)
			response := protocol.ErrorValue(err.Error())
			conn.Write([]byte(response.Marshal()))
			continue
		}

		// Log the command
		cmdStr := commandToString(command, args)
		LogCommand(clientAddr, cmdStr)

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

// commandToString converts a command and its args to a loggable string
func commandToString(command string, args []protocol.ORSPValue) string {
	var sb strings.Builder
	sb.WriteString(command)

	for _, arg := range args {
		switch v := arg.(type) {
		case protocol.BulkStringValue:
			sb.WriteString(" \"")
			sb.WriteString(string(v))
			sb.WriteString("\"")
		default:
			sb.WriteString(" ")
			sb.WriteString(fmt.Sprintf("%v", v))
		}
	}

	return sb.String()
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
