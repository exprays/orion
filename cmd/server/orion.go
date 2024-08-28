package main

import (
	"flag"
	"fmt"
	"os"

	"orion/src/aof"
	"orion/src/server"
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
		// Load AOF to restore state
		// fmt.Println("Loading AOF data...")
		// err = aof.LoadAOF(func(command protocol.ArrayValue) error {
		// 	// Convert ArrayValue to string command or handle it as needed
		// 	commandStr := command.Marshal()

		// 	// Parse the string command
		// 	args := parseStringCommand(commandStr)

		// 	// Handle the command
		// 	response := server.HandleCommand(args)
		// 	if _, ok := response.(protocol.ErrorValue); ok {
		// 		// return fmt.Errorf("error from server while handling command: %s", commandStr)
		// 	}
		// 	return nil
		// })

		// if err != nil {
		// 	fmt.Printf("Error loading AOF: %v\n", err)
		// }
		// if err != nil {
		// 	fmt.Println("Error loading AOF file:", err)
		// 	os.Exit(1)
		// }

		server.StartServer(*port)
	} else {
		fmt.Println("Unknown mode. Use `server`.")
	}
}

// parseStringCommand converts a string command to []protocol.ORSPValue
// func parseStringCommand(command string) []protocol.ORSPValue {
// 	parts := strings.Fields(command)
// 	args := make([]protocol.ORSPValue, len(parts))
// 	for i, part := range parts {
// 		args[i] = protocol.BulkStringValue(part)
// 	}
// 	return args
// }
