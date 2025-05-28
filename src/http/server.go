package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"orion/src/data"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type HTTPServer struct {
	port     int
	router   *mux.Router
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
	store    *data.DataStore // Add reference to the store
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

func NewHTTPServer(port int, store *data.DataStore) *HTTPServer {
	server := &HTTPServer{
		port:  port,
		store: store, // Initialize with store
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

func (s *HTTPServer) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	uptime := time.Since(startTime)

	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Get total keys from store
	totalKeys := s.store.Size()

	stats := ServerStats{
		UptimeInSeconds: int(uptime.Seconds()),
		UptimeInDays:    int(uptime.Hours() / 24),
		UsedMemory:      int64(m.Alloc),
		UsedMemoryHuman: formatBytes(int64(m.Alloc)),
		KeyspaceInfo:    fmt.Sprintf("db0:keys=%d", totalKeys),
		TotalKeys:       totalKeys,
		Connections:     len(s.clients),
		Version:         "Orion 1.0.0",
		Port:            s.port,
	}

	response := APIResponse{
		Success: true,
		Data:    stats,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *HTTPServer) handleKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	pattern := r.URL.Query().Get("pattern")
	if pattern == "" {
		pattern = "*"
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Get keys from store
	keys := s.store.GetAllData()

	var keyInfos []KeyInfo
	count := 0

	for _, key := range keys {
		if count >= limit {
			break
		}

		// Simple pattern matching (for now, just support * for all)
		if pattern != "*" && !strings.Contains(key, pattern) {
			continue
		}

		// Get key info
		value, exists := s.store.Get(key)
		if !exists {
			continue
		}

		keyInfo := KeyInfo{
			Key:   key,
			Type:  "string", // For now, assume all are strings
			Value: value,
			TTL:   -1, // No TTL support yet
			Size:  len(value),
		}

		keyInfos = append(keyInfos, keyInfo)
		count++
	}

	response := APIResponse{
		Success: true,
		Data:    keyInfos,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *HTTPServer) handleKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	keyName := vars["name"]

	if keyName == "" {
		response := APIResponse{
			Success: false,
			Error:   "Key name is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	switch r.Method {
	case "GET":
		value, exists := s.store.Get(keyName)
		if !exists {
			response := APIResponse{
				Success: false,
				Error:   "Key not found",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		keyInfo := KeyInfo{
			Key:   keyName,
			Type:  "string",
			Value: value,
			TTL:   -1,
			Size:  len(value),
		}

		response := APIResponse{
			Success: true,
			Data:    keyInfo,
		}
		json.NewEncoder(w).Encode(response)

	case "DELETE":
		existed := s.store.Delete(keyName)
		if !existed {
			response := APIResponse{
				Success: false,
				Error:   "Key not found",
			}
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(response)
			return
		}

		response := APIResponse{
			Success: true,
			Data:    map[string]string{"message": "Key deleted successfully"},
		}
		json.NewEncoder(w).Encode(response)
	}
}

func (s *HTTPServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	s.clients[conn] = true
	log.Printf("WebSocket client connected. Total clients: %d", len(s.clients))

	// Handle incoming messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(s.clients, conn)
			log.Printf("WebSocket client disconnected. Total clients: %d", len(s.clients))
			break
		}
	}
}
