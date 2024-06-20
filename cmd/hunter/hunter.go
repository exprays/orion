package main

import (
	"fmt"
	"os"

	"orion/src/hunter" // Import the hunter package
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: hunter connect")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "connect":
		hunter.Connect()
	default:
		fmt.Println("Unknown command. Available commands: connect")
	}
}
