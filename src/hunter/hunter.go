package hunter

import (
	"bufio"
	"fmt"
	"net"
	"orion/src/protocol"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/chzyer/readline"
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
	color.Magenta("Made with love by The Star Society")

	mode := promptInput("Select mode (dev/custom): ", color.FgGreen)
	mode = strings.ToLower(mode)

	var serverIP, serverPort string

	//dev mode is the default mode of hunter which is used for local development
	//custom mode is used for connecting to a custom server

	if mode == "dev" {
		serverIP = "127.0.0.1"
		serverPort = "6379"
		color.Green("Dev mode selected. Using IP: %s and Port: %s", serverIP, serverPort)
	} else if mode == "custom" {
		serverIP = promptInput("Enter server IP: ", color.FgGreen)
		serverPort = promptInput("Enter server port: ", color.FgGreen)
	} else {
		color.Red("Invalid mode selected. Exiting...")
		return
	}

	serverAddr := fmt.Sprintf("%s:%s", serverIP, serverPort)

	showLoader()

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		color.Red("Error connecting to server: %v", err)
		return
	}
	defer conn.Close()

	color.Green("Connected to server at %s", serverAddr)

	// Initialize command history
	history := NewCommandHistory()

	// Setup readline with proper configuration
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          color.GreenString("%s> ", serverAddr),
		HistoryFile:     "", // We're managing history ourselves
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		// Tab completion function
		AutoComplete: newAutoCompleter(),
	})
	if err != nil {
		color.Red("Error initializing readline: %v", err)
		return
	}
	defer rl.Close()

	// Setup signal handler to catch ctrl+c
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nCtrl+C detected. Exiting...")
		os.Exit(0)
	}()

	for {
		input, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Handle local commands
		if handleLocalCommand(input, history) {
			continue
		}

		// Add to history
		history.Add(input)

		// Convert input to ORSP array
		args := parseInput(input)
		orspArray := make(protocol.ArrayValue, len(args))
		for i, arg := range args {
			orspArray[i] = protocol.BulkStringValue(arg)
		}

		// Marshal and send the ORSP array
		_, err = conn.Write([]byte(orspArray.Marshal()))
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

// handleLocalCommand processes commands that are handled locally by the CLI
func handleLocalCommand(input string, history *CommandHistory) bool {
	cmd := strings.ToUpper(strings.Fields(input)[0])

	switch cmd {
	case "CLEAR":
		clearScreen()
		return true
	case "HISTORY":
		showHistory(history)
		return true
	case "HELP":
		showHelp()
		return true
	case "EXIT", "QUIT":
		fmt.Println("Goodbye!")
		os.Exit(0)
	}
	return false
}

func showHistory(history *CommandHistory) {
	entries := history.List(20) // Show last 20 commands
	if len(entries) == 0 {
		color.Yellow("No command history")
		return
	}

	color.Cyan("Command History:")
	for i, entry := range entries {
		color.White("%d. %s [%s]", i+1, entry.Command, entry.Timestamp.Format("2006-01-02 15:04:05"))
	}
}

func showHelp() {
	color.Cyan("Hunter CLI Help:")
	help := []struct {
		cmd  string
		desc string
	}{
		{"CLEAR", "Clear the screen"},
		{"HISTORY", "Show command history"},
		{"HELP", "Show this help message"},
		{"EXIT/QUIT", "Exit the CLI"},
	}

	for _, h := range help {
		fmt.Printf("%s%s%s\n", color.GreenString(h.cmd), strings.Repeat(" ", 15-len(h.cmd)), h.desc)
	}

	color.Yellow("\nServer Commands:")
	color.Yellow("For a list of all available server commands, type: INFO")
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

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()

	// Re-print the ASCII art and welcome message
	color.Cyan(asciiArt)
	color.Yellow("Welcome to Hunter CLI 1.0!")
	color.Yellow("Read more about hunter on https://orion.thestarsociety.tech/docs/packages/hunter")
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
