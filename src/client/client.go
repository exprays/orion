package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
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
