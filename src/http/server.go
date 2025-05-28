package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type HTTPServer struct {
	port     int
	router   *mux.Router
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ServerStats struct {
	UptimeInSeconds int    `json:"uptime_in_seconds"`
	UptimeInDays    int    `json:"uptime_in_days"`
	UsedMemory      int64  `json:"used_memory"`
	UsedMemoryHuman string `json:"used_memory_human"`
	KeyspaceInfo    string `json:"keyspace_info"`
	TotalKeys       int    `json:"total_keys"`
	Connections     int    `json:"connections"`
	Version         string `json:"version"`
	Port            int    `json:"port"`
}

type KeyInfo struct {
	Key     string            `json:"key"`
	Type    string            `json:"type"`
	Value   string            `json:"value,omitempty"`
	Members []string          `json:"members,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
	TTL     int               `json:"ttl"`
	Size    int               `json:"size"`
}

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResult struct {
	Command   string `json:"command"`
	Result    string `json:"result"`
	Timestamp int64  `json:"timestamp"`
	Success   bool   `json:"success"`
}

var (
	startTime = time.Now()
)

func NewHTTPServer(port int) *HTTPServer {
	server := &HTTPServer{
		port: port,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow connections from localhost for development
				origin := r.Header.Get("Origin")
				return strings.Contains(origin, "localhost") ||
					strings.Contains(origin, "127.0.0.1") ||
					origin == "" // Allow same-origin requests
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}

	server.setupRoutes()
	return server
}

func (s *HTTPServer) setupRoutes() {
	s.router = mux.NewRouter()

	// Add CORS middleware
	s.router.Use(s.corsMiddleware)
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.securityMiddleware)

	// API routes
	api := s.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/stats", s.handleStats).Methods("GET", "OPTIONS")
	api.HandleFunc("/keys", s.handleKeys).Methods("GET", "OPTIONS")
	api.HandleFunc("/key/{name}", s.handleKey).Methods("GET", "DELETE", "OPTIONS")
	api.HandleFunc("/command", s.handleCommand).Methods("POST", "OPTIONS")

	// WebSocket route
	s.router.HandleFunc("/ws", s.handleWebSocket)

	// Health check
	s.router.HandleFunc("/health", s.handleHealth).Methods("GET")
}

func (s *HTTPServer) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("HTTP server starting on port %d", s.port)
	log.Printf("Dashboard available at: http://localhost:%d", s.port)

	return http.ListenAndServe(addr, s.router)
}
