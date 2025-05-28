package http

import (
	"log"
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
