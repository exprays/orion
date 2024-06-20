package main

import (
	"flag"
	"fmt"
	"os"

	"orion/src/aof"    // Import the AOF package
	"orion/src/server" // Import the server package
)

func main() {
	mode := flag.String("mode", "server", "start in `server` mode")
	port := flag.String("port", "6379", "port to run the server on")
	flag.Parse()

	if *mode == "server" {
		// Initialize AOF
		err := aof.InitAOF()
		if err != nil {
			fmt.Println("Error initializing AOF:", err)
			os.Exit(1)
		}

		// Load AOF file and handle errors
		err = aof.LoadAOF(func(command string) error {
			response := server.HandleCommand(command)
			if response == "ERR" {
				return fmt.Errorf("error from server while handling command: %s", command)
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error loading AOF file:", err)
			os.Exit(1)
		}

		server.StartServer(*port)
	} else {
		fmt.Println("Unknown mode. Use `server`.")
	}
}
