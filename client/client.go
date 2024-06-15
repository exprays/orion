package client

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
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
