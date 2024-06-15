package main

import (
    "flag"
    "fmt"
    "os"

    "orion/client"
    "orion/server"
)

func main() {
    mode := flag.String("mode", "server", "start in `server` or `client` mode")
    port := flag.String("port", "6379", "port to run the server on or to connect to")
    flag.Parse()

    if *mode == "server" {
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
