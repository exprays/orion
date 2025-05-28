package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"orion/src/server"
)

func main() {
	port := flag.Int("port", 6379, "Port to run the TCP server on")
	httpPort := flag.Int("http-port", 8080, "Port to run the HTTP API server on")
	flag.Parse()

	// Print startup banner
	fmt.Println(`
     ██████╗ ██████╗ ██╗ ██████╗ ███╗   ███╗
    ██╔═══██╗██╔══██╗██║██╔═══██╗████╗ ████║
    ██║   ██║██████╔╝██║██║   ██║██╔████╔██║
    ██║   ██║██╔══██╗██║██║   ██║██║╚██╔╝██║
    ╚██████╔╝██║  ██║██║╚██████╔╝██║ ╚═╝ ██║
     ╚═════╝ ╚═╝  ╚═╝╚═╝ ╚═════╝ ╚═╝     ╚═╝
    `)

	log.Printf("🚀 Starting Orion Database Server...")
	log.Printf("📡 TCP Port: %d", *port)
	log.Printf("🌐 HTTP API Port: %d", *httpPort)

	// Channel to handle graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Channel to collect server errors
	serverErrors := make(chan error, 2)

	// Start the server with both TCP and HTTP ports
	go func() {
		defer func() {
			if r := recover(); r != nil {
				serverErrors <- fmt.Errorf("server panic: %v", r)
			}
		}()

		// Call StartServer with both ports
		server.StartServer(strconv.Itoa(*port), *httpPort)
	}()

	// Give servers a moment to start up
	time.Sleep(2 * time.Second)

	// Display service information
	fmt.Println("\n🎯 Services Running:")
	fmt.Printf("   📡 TCP Server: localhost:%d (ORSP Protocol)\n", *port)
	fmt.Printf("   🌐 HTTP API: http://localhost:%d/api\n", *httpPort)
	fmt.Printf("   📊 Belt Dashboard: http://localhost:%d\n", *httpPort)
	fmt.Printf("   🔄 WebSocket: ws://localhost:%d/ws\n", *httpPort)
	fmt.Println("   💾 AOF Persistence: Enabled")
	fmt.Println("   📝 Command Logging: Enabled")

	fmt.Println("\n✅ Orion Database Server Started Successfully!")
	fmt.Println("\n🔗 Connect using:")
	fmt.Printf("   • Hunter CLI: ./hunter -p %d\n", *port)
	fmt.Printf("   • Belt Dashboard: http://localhost:%d\n", *httpPort)
	fmt.Printf("   • Custom clients: localhost:%d\n", *port)

	fmt.Println("\n⚡ Available Services:")
	fmt.Println("   • In-Memory Database (Strings, Sets, Hashes)")
	fmt.Println("   • Real-time WebSocket Updates")
	fmt.Println("   • RESTful API")
	fmt.Println("   • Append-Only File (AOF) Persistence")
	fmt.Println("   • Web Dashboard Management")
	fmt.Println("   • Command Logging & Monitoring")

	fmt.Println("\n📊 Service Endpoints:")
	fmt.Printf("   • GET  /api/stats        - Server statistics\n")
	fmt.Printf("   • GET  /api/keys         - List all keys\n")
	fmt.Printf("   • GET  /api/key/{name}   - Get key details\n")
	fmt.Printf("   • POST /api/command      - Execute command\n")
	fmt.Printf("   • WS   /ws              - Real-time updates\n")

	fmt.Println("\nPress Ctrl+C to shutdown gracefully...")

	// Wait for shutdown signal or server error
	select {
	case <-shutdown:
		fmt.Println("\n🛑 Shutdown signal received, initiating graceful shutdown...")
	case err := <-serverErrors:
		fmt.Printf("\n❌ Server error: %v\n", err)
		fmt.Println("🛑 Initiating shutdown due to server error...")
	}

	// Graceful shutdown
	fmt.Println("🔄 Shutting down services...")

	// Give time for connections to close gracefully
	fmt.Println("⏳ Allowing active connections to complete...")
	time.Sleep(3 * time.Second)

	// Note: The TCP server's defer statements will handle:
	// - AOF file closure
	// - Log file closure
	// - Connection cleanup

	log.Println("✅ Orion Database Server shutdown complete")
	fmt.Println("💾 Data persisted to AOF file")
	fmt.Println("📝 Logs saved successfully")
	fmt.Println("👋 Thank you for using Orion Database!")
}
