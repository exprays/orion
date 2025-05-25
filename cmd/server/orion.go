package main

import (
	"flag"
	"fmt"

	"orion/src/server"
)

func main() {
	mode := flag.String("mode", "server", "start in `server` mode")
	port := flag.String("port", "6379", "port to run the server on")
	flag.Parse()

	if *mode == "server" {
		fmt.Println("Starting Orion server...")
		server.StartServer(*port)
	} else {
		fmt.Println("Unknown mode. Use `server`.")
	}
}
