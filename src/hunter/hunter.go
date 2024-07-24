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
	"time"

	"github.com/fatih/color"
)

var asciiArt = `
  _______   ________  ________  ________  ________  ________ 
  /    /  \\/    /   \/    /   \/        \/        \/        \
 /        //         /         /        _/         /         /
/         /         /         //       //        _/        _/ 
\___/____/\________/\__/_____/ \______/ \________/\____/___/  
`

// Connect initializes the CLI client and connects to the server.
func Connect() {
	color.Cyan(asciiArt)
	color.Yellow("Welcome to Hunter CLI 1.0!")
	color.Yellow("Read more about hunter on https://orion.thestarsociety.tech/docs/packages/hunter")

	serverIP := promptInput("Enter server IP: ", color.FgGreen)
	serverPort := promptInput("Enter server port: ", color.FgGreen)
	serverAddr := fmt.Sprintf("%s:%s", serverIP, serverPort)

	showLoader()

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		color.Red("Error connecting to server: %v", err)
		return
	}
	defer conn.Close()

	color.Green("Connected to server at %s", serverAddr)

	// Setup signal handler to catch ctrl+c
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nCtrl+C detected. Exiting...")
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		prompt := color.WhiteString("%s> ", serverAddr)
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		// Convert input to ORSP array
		args := parseInput(input)
		orspArray := make(protocol.ArrayValue, len(args))
		for i, arg := range args {
			orspArray[i] = protocol.BulkStringValue(arg)
		}

		// Marshal and send the ORSP array
		_, err := conn.Write([]byte(orspArray.Marshal()))
		if err != nil {
			color.Red("Error sending command: %v", err)
			continue
		}

		// Read and unmarshal the response
		respReader := bufio.NewReader(conn)
		response, err := protocol.Unmarshal(respReader)
		if err != nil {
			color.Red("Error reading response: %v", err)
			continue
		}

		// Print the response
		printResponse(response)
	}
}

func promptInput(prompt string, textColor color.Attribute) string {
	fmt.Print(color.New(textColor).Sprint(prompt))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func showLoader() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	for i := 0; i < 20; i++ {
		frame := frames[i%len(frames)]
		fmt.Printf("\r%s Connecting...", color.CyanString(frame))
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println()
}

func printResponse(response protocol.ORSPValue) {
	switch v := response.(type) {
	case protocol.SimpleStringValue:
		color.Green("%s", string(v))
	case protocol.ErrorValue:
		color.Red("Error: %s", string(v))
	case protocol.IntegerValue:
		color.Blue("%d", int64(v))
	case protocol.BulkStringValue:
		color.Cyan("%s", string(v))
	case protocol.ArrayValue:
		for _, item := range v {
			printResponse(item)
		}
	case protocol.NullValue:
		color.Magenta("(nil)")
	case protocol.BooleanValue:
		color.Yellow("%v", bool(v))
	case protocol.DoubleValue:
		color.Blue("%f", float64(v))
	case *protocol.BigNumberValue:
		color.Blue("%s", v.String())
	case protocol.BulkErrorValue:
		color.Red("Error (%s): %s", v.Code, v.Message)
	case protocol.VerbatimStringValue:
		color.Cyan("%s:%s", v.Format, v.Value)
	case protocol.MapValue:
		for key, value := range v {
			fmt.Printf("%s: ", color.HiYellowString(key))
			printResponse(value)
		}
	case protocol.SetValue:
		for _, item := range v {
			printResponse(item)
		}
	case protocol.PushValue:
		color.HiMagenta("Push (%s):", v.Kind)
		for _, item := range v.Data {
			printResponse(item)
		}
	default:
		color.HiRed("Unknown type: %T", v)
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
			if !inQuotes && currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
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
