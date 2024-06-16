package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
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

		respCommand := formatRespCommand(input)
		fmt.Printf("DEBUG: Sending RESP command: %s\n", respCommand) // Debugging line
		_, err := conn.Write([]byte(respCommand))
		if err != nil {
			fmt.Println("Error sending command:", err)
			continue
		}

		response, err := readResp(conn)
		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		fmt.Println(response)
	}
}

// formatRespCommand converts a user input command into RESP format
func formatRespCommand(command string) string {
	parts := parseArguments(command)
	if len(parts) == 0 {
		return ""
	}

	// Format as RESP array
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("*%d\r\n", len(parts)))
	for _, part := range parts {
		builder.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(part), part))
	}
	return builder.String()
}

// parseArguments splits the command into arguments while handling quoted strings
func parseArguments(input string) []string {
	// Regex to match words or quoted strings
	re := regexp.MustCompile(`"([^"]*)"|(\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	var args []string
	for _, match := range matches {
		if match[1] != "" {
			// Quoted string
			args = append(args, match[1])
		} else {
			// Regular word
			args = append(args, match[2])
		}
	}
	return args
}

// readResp reads and parses RESP data from the server
func readResp(conn net.Conn) (string, error) {
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
			elem, err := readResp(conn)
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
			elem, err := readResp(conn)
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
			elem, err := readResp(conn)
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
			elem, err := readResp(conn)
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
