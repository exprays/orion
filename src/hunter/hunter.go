package hunter

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Connect initializes the CLI client and connects to the server.
func Connect() {
	var serverAddr string

	// Prompt user for server address
	fmt.Print("Welcome to Hunter CLI!\n")
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

		fmt.Fprintf(conn, input+"\n")
		response, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(response)
	}
}
