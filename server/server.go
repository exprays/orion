package server

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

// StartServer initializes the TCP server
func StartServer(port string) {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Server is running and listening on port", port)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    for {
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Error reading input:", err)
            return
        }
        input = strings.TrimSpace(input)
        response := HandleCommand(input)
        conn.Write([]byte(response + "\n"))
    }
}
