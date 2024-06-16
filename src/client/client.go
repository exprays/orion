package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// StartClient initializes the CLI client
func StartClient(serverAddr string) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server at", serverAddr)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		thunderCommand := formatThunderCommand(input)
		_, err := conn.Write([]byte(thunderCommand))
		if err != nil {
			fmt.Println("Error sending command:", err)
			continue
		}

		response, err := readThunder(conn)
		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		fmt.Println(response)
	}
}

// formatThunderCommand converts a user input command into THUNDER format
func formatThunderCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}

	// Format as Thunder array
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("*%d\r\n", len(parts)))
	for _, part := range parts {
		// Trim quotes if the argument is enclosed in quotes
		part = strings.Trim(part, `"`)
		builder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(part), part))
	}
	return builder.String()
}

// readThunder reads and parses THUNDER data from the server
func readThunder(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)

	if len(line) == 0 {
		return "", nil
	}

	switch line[0] {
	case '+': // Simple String
		return line[1:], nil
	case '-': // Error
		return "", fmt.Errorf(line[1:])
	case ':': // Integer
		return line[1:], nil
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
	case '*': // Array
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", err
		}
		var result []string
		for i := 0; i < count; i++ {
			elem, err := readThunder(conn)
			if err != nil {
				return "", err
			}
			result = append(result, elem)
		}
		return strings.Join(result, " "), nil
	case '_': // Null
		return "", nil
	case '#': // Boolean
		if line[1] == 't' {
			return "true", nil
		} else {
			return "false", nil
		}
	case ',': // Double
		return line[1:], nil
	case '(': // Big Number
		return line[1:], nil
	case '=': // Verbatim String
		return line[1:], nil
	case '%': // Map
		// For simplicity, we assume a flat key-value map
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", err
		}
		var result []string
		for i := 0; i < count*2; i++ { // Key and value pairs
			elem, err := readThunder(conn)
			if err != nil {
				return "", err
			}
			result = append(result, elem)
		}
		return fmt.Sprintf("{%s}", strings.Join(result, " ")), nil
	case '~': // Set
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", err
		}
		var result []string
		for i := 0; i < count; i++ {
			elem, err := readThunder(conn)
			if err != nil {
				return "", err
			}
			result = append(result, elem)
		}
		return fmt.Sprintf("[%s]", strings.Join(result, " ")), nil
	case '>': // Push
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", err
		}
		var result []string
		for i := 0; i < count; i++ {
			elem, err := readThunder(conn)
			if err != nil {
				return "", err
			}
			result = append(result, elem)
		}
		return strings.Join(result, " "), nil
	default:
		return "", fmt.Errorf("invalid thunder format: %s", line)
	}
}
