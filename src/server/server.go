package server

import (
	"bufio"
	"fmt"
	"net"
	"orion/src/aof"
	"strings"
	"unicode"
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
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		input = strings.TrimSpace(input)

		command, args, err := parseCommand(input)
		if err != nil {
			conn.Write([]byte("ERROR: " + err.Error() + "\n"))
			continue
		}

		fullCommand := command + " " + strings.Join(args, " ")
		response := HandleCommand(fullCommand)
		conn.Write([]byte(response + "\n"))

		if command != "" {
			aof.AppendCommand(fullCommand)
		}
	}
}

func parseCommand(input string) (string, []string, error) {
	var args []string
	var currentArg strings.Builder
	inQuotes := false
	escaped := false

	for _, char := range input {
		if escaped {
			currentArg.WriteRune(char)
			escaped = false
			continue
		}
		if char == '\\' {
			escaped = true
			continue
		}
		if char == '"' {
			inQuotes = !inQuotes
			continue
		}
		if unicode.IsSpace(char) && !inQuotes {
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
			continue
		}
		currentArg.WriteRune(char)
	}
	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	if inQuotes {
		return "", nil, fmt.Errorf("unmatched quote in input")
	}

	if len(args) == 0 {
		return "", nil, fmt.Errorf("no command provided")
	}

	return strings.ToUpper(args[0]), args[1:], nil

}
