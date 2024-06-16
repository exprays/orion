// server/server.go
package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// StartServer initializes the TCP server with THUNDER protocol
func StartServer(port string) {
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
		input, err := readThunder(reader)
		if err != nil {
			fmt.Println("Error reading input:", err)
			conn.Write([]byte("-ERROR reading input\r\n"))
			return
		}
		response := HandleCommand(input)
		conn.Write([]byte(response))
	}
}

// readThunder reads and parses THUNDER data from the connection
func readThunder(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return "", nil
	}

	switch line[0] {
	case '*': // Array
		numArgs, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", fmt.Errorf("invalid number of arguments: %s", line)
		}
		args := make([]string, numArgs)
		for i := 0; i < numArgs; i++ {
			arg, err := readThunder(reader)
			if err != nil {
				return "", err
			}
			args[i] = arg
		}
		return strings.Join(args, " "), nil
	case '$': // Bulk String
		length, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", err
		}
		if length == -1 {
			return "", nil // Null bulk string
		}
		data := make([]byte, length+2) // Include \r\n
		_, err = reader.Read(data)
		if err != nil {
			return "", err
		}
		return string(data[:length]), nil
	case '+', '-', ':': // Simple String, Error, or Integer
		return line[1:], nil
	case '_': // Null
		return "", nil
	case '#': // Boolean
		return line[1:], nil
	case ',': // Double
		return line[1:], nil
	case '(': // Big Number
		return line[1:], nil
	case '=': // Verbatim String
		return line[1:], nil
	case '%', '~': // Maps or Sets
		return line[1:], nil
	case '>': // Pushes
		return line[1:], nil
	default:
		return "", fmt.Errorf("invalid thunder format: %s", line)
	}
}
