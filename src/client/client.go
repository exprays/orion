package hunter

import (
	"bufio"
	"fmt"
	"net"
	"orion/src/protocol"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Connect initializes the CLI client and connects to the server.
func Connect() {
	var serverAddr string

	// Prompt user for server address
	fmt.Print("Welcome to Hunter CLI 1.0!\n")
	fmt.Print("Read more about hunter on https://orion.thestarsociety.tech/docs/packages/hunter\n")
	fmt.Print("Enter server address (IP:Port): ")
	reader := bufio.NewReader(os.Stdin)
	serverAddr, _ = reader.ReadString('\n')
	serverAddr = strings.TrimSpace(serverAddr)

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server at", serverAddr)

	// Setup signal handler to catch ctrl+c
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nCtrl+C detected. Exiting...")
		os.Exit(0)
	}()

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		// Parse the input, respecting quoted strings
		args := parseInput(input)

		// Convert input to ORSP array
		orspArray := make(protocol.ArrayValue, len(args))
		for i, arg := range args {
			orspArray[i] = protocol.BulkStringValue(arg)
		}

		// Marshal and send the ORSP array
		_, err := conn.Write([]byte(orspArray.Marshal()))
		if err != nil {
			fmt.Println("Error sending command:", err)
			continue
		}

		// Read and unmarshal the response
		respReader := bufio.NewReader(conn)
		response, err := protocol.Unmarshal(respReader)
		if err != nil {
			fmt.Println("Error reading response:", err)
			continue
		}

		// Print the response
		printResponse(response)
	}
}

func parseInput(input string) []string {
	var args []string
	var currentArg strings.Builder
	inQuotes := false

	for _, char := range input {
		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if !inQuotes {
				if currentArg.Len() > 0 {
					args = append(args, currentArg.String())
					currentArg.Reset()
				}
			} else {
				currentArg.WriteRune(char)
			}
		default:
			currentArg.WriteRune(char)
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	return args
}

func printResponse(response protocol.ORSPValue) {
	switch v := response.(type) {
	case protocol.SimpleStringValue:
		fmt.Println(string(v))
	case protocol.ErrorValue:
		fmt.Println("Error:", string(v))
	case protocol.IntegerValue:
		fmt.Println(int64(v))
	case protocol.BulkStringValue:
		fmt.Printf("\"%s\"\n", string(v))
	case protocol.ArrayValue:
		for i, item := range v {
			if i > 0 {
				fmt.Print(" ")
			}
			printResponse(item)
		}
		fmt.Println()
	case protocol.NullValue:
		fmt.Println("(nil)")
	case protocol.BooleanValue:
		fmt.Println(bool(v))
	case protocol.DoubleValue:
		fmt.Println(float64(v))
	case *protocol.BigNumberValue:
		fmt.Println(v.String())
	case protocol.BulkErrorValue:
		fmt.Printf("Error (%s): %s\n", v.Code, v.Message)
	case protocol.VerbatimStringValue:
		fmt.Printf("%s:%s\n", v.Format, v.Value)
	case protocol.MapValue:
		for key, value := range v {
			fmt.Printf("%s: ", key)
			printResponse(value)
		}
	case protocol.SetValue:
		for _, item := range v {
			printResponse(item)
		}
	case protocol.PushValue:
		fmt.Printf("Push (%s):\n", v.Kind)
		for _, item := range v.Data {
			printResponse(item)
		}
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}
