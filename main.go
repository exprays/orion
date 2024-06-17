package main

import (
	"flag"
	"fmt"
	"os"

	"orion/src/aof" // Import the AOF package
	"orion/src/client"
	"orion/src/server"
)

func main() {
	mode := flag.String("mode", "server", "start in `server` or `client` mode")
	port := flag.String("port", "6379", "port to run the server on or to connect to")
	flag.Parse()

	if *mode == "server" {
		// Initialize AOF
		err := aof.InitAOF()
		if err != nil {
			fmt.Println("Error initializing AOF:", err)
			os.Exit(1)
		}

		// **Fix 1: Wrapper function for error handling**
		err = aof.LoadAOF(func(command string) error {
			response := server.HandleCommand(command)
			if response == "ERR" {
				// Handle error case based on the response
				return fmt.Errorf("error from server while handling command: %s", command)
			}
			return nil // Or return the error from server.HandleCommand
		})
		if err != nil {
			fmt.Println("Error loading AOF file:", err)
			os.Exit(1)
		}

		server.StartServer(*port)
	} else if *mode == "client" {
		if flag.NArg() == 0 {
			fmt.Println("Usage: go run main.go --mode=client <server-address>")
			os.Exit(1)
		}
		serverAddr := flag.Arg(0)
		client.StartClient(serverAddr)
	} else {
		fmt.Println("Unknown mode. Use `server` or `client`.")
	}
}
