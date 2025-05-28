package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type KeyUpdateMessage struct {
	Action string `json:"action"`
	Key    string `json:"key"`
}

type CommandExecutionMessage struct {
	Command string `json:"command"`
	Result  string `json:"result"`
	Success bool   `json:"success"`
}

func (s *HTTPServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Add client to the list
	s.clients[conn] = true
	log.Printf("WebSocket client connected. Total clients: %d", len(s.clients))

	// Send initial connection message
	s.sendToClient(conn, WebSocketMessage{
		Type:      "connection",
		Data:      map[string]string{"status": "connected"},
		Timestamp: time.Now().Unix(),
	})

	// Handle incoming messages
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Handle different message types
		msgType, ok := msg["type"].(string)
		if !ok {
			continue
		}

		switch msgType {
		case "ping":
			s.sendToClient(conn, WebSocketMessage{
				Type:      "pong",
				Data:      map[string]string{"status": "alive"},
				Timestamp: time.Now().Unix(),
			})
		case "subscribe":
			// Handle subscription to specific events
			s.handleSubscription(conn, msg)
		}
	}

	// Remove client from the list
	delete(s.clients, conn)
	log.Printf("WebSocket client disconnected. Total clients: %d", len(s.clients))
}

func (s *HTTPServer) handleSubscription(conn *websocket.Conn, msg map[string]interface{}) {
	// Handle different subscription types
	events, ok := msg["events"].([]interface{})
	if !ok {
		return
	}

	// For now, we'll just acknowledge the subscription
	s.sendToClient(conn, WebSocketMessage{
		Type: "subscription_ack",
		Data: map[string]interface{}{
			"events": events,
		},
		Timestamp: time.Now().Unix(),
	})
}

func (s *HTTPServer) sendToClient(conn *websocket.Conn, message WebSocketMessage) {
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err := conn.WriteJSON(message); err != nil {
		log.Printf("WebSocket write error: %v", err)
		conn.Close()
		delete(s.clients, conn)
	}
}

func (s *HTTPServer) broadcastToClients(message WebSocketMessage) {
	for conn := range s.clients {
		go s.sendToClient(conn, message)
	}
}

func (s *HTTPServer) broadcastKeyUpdate(action, key string) {
	message := WebSocketMessage{
		Type: "key_update",
		Data: KeyUpdateMessage{
			Action: action,
			Key:    key,
		},
		Timestamp: time.Now().Unix(),
	}
	s.broadcastToClients(message)
}

func (s *HTTPServer) broadcastCommandExecution(command, result string, success bool) {
	message := WebSocketMessage{
		Type: "command_execution",
		Data: CommandExecutionMessage{
			Command: command,
			Result:  result,
			Success: success,
		},
		Timestamp: time.Now().Unix(),
	}
	s.broadcastToClients(message)
}
