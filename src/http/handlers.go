package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"orion/src/commands"
	"orion/src/core"
	"orion/src/data"
	"orion/src/protocol"

	"github.com/gorilla/mux"
)

func (s *HTTPServer) handleStats(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	totalKeys := core.GetKeyCount()

	stats := ServerStats{
		UptimeInSeconds: int(uptime.Seconds()),
		UptimeInDays:    int(uptime.Hours() / 24),
		UsedMemory:      int64(memStats.Alloc),
		UsedMemoryHuman: formatBytes(int64(memStats.Alloc)),
		KeyspaceInfo:    fmt.Sprintf("db0:keys=%d", totalKeys),
		TotalKeys:       totalKeys,
		Connections:     len(s.clients),
		Version:         "v0.1.0",
		Port:            s.port,
	}

	s.writeJSON(w, APIResponse{
		Success: true,
		Data:    stats,
	})
}

func (s *HTTPServer) handleKeys(w http.ResponseWriter, r *http.Request) {
	pattern := r.URL.Query().Get("pattern")
	if pattern == "" {
		pattern = "*"
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	keys := core.GetKeys(pattern, limit)
	keyInfos := make([]KeyInfo, 0, len(keys))

	for _, key := range keys {
		keyInfo := s.getKeyInfo(key)
		if keyInfo != nil {
			keyInfos = append(keyInfos, *keyInfo)
		}
	}

	s.writeJSON(w, APIResponse{
		Success: true,
		Data:    keyInfos,
	})
}

func (s *HTTPServer) handleKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keyName := vars["name"]

	if keyName == "" {
		s.writeError(w, http.StatusBadRequest, "Key name is required")
		return
	}

	switch r.Method {
	case "GET":
		keyInfo := s.getKeyInfo(keyName)
		if keyInfo == nil {
			s.writeError(w, http.StatusNotFound, "Key not found")
			return
		}

		s.writeJSON(w, APIResponse{
			Success: true,
			Data:    keyInfo,
		})

	case "DELETE":
		// Execute DEL command using the data store directly
		data.Store.Del(keyName)

		// Notify WebSocket clients
		s.broadcastKeyUpdate("delete", keyName)

		s.writeJSON(w, APIResponse{
			Success: true,
			Data:    true,
		})
	}
}

func (s *HTTPServer) handleCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.Command == "" {
		s.writeError(w, http.StatusBadRequest, "Command is required")
		return
	}

	// Parse and execute command
	parts := strings.Fields(req.Command)
	if len(parts) == 0 {
		s.writeError(w, http.StatusBadRequest, "Invalid command")
		return
	}

	// Convert to ORSP values
	orspValues := make([]protocol.ORSPValue, len(parts))
	for i, part := range parts {
		orspValues[i] = protocol.BulkStringValue(part)
	}

	// Execute command using existing command handlers
	var result protocol.ORSPValue
	var success bool = true

	commandName := strings.ToUpper(parts[0])
	args := orspValues[1:]

	switch commandName {
	// String Commands
	case "GET":
		result = commands.HandleGet(args)
	case "SET":
		result = commands.HandleSet(args)
	case "APPEND":
		result = commands.HandleAppend(args)
	case "GETDEL":
		result = commands.HandleGetDel(args)
	case "GETEX":
		result = commands.HandleGetEx(args)
	case "GETSET":
		result = commands.HandleGetSet(args)
	case "GETRANGE":
		result = commands.HandleGetRange(args)
	case "INCR":
		result = commands.HandleIncr(args)
	case "INCRBY":
		result = commands.HandleIncrBy(args)
	case "INCRBYFLOAT":
		result = commands.HandleIncrByFloat(args)
	case "LCS":
		result = commands.HandleLCS(args)

	// TTL Commands
	case "TTL":
		result = commands.HandleTTL(args)

	// Set Commands
	case "SADD":
		result = commands.HandleSAdd(args)
	case "SCARD":
		result = commands.HandleSCard(args)
	case "SMEMBERS":
		result = commands.HandleSMembers(args)
	case "SISMEMBER":
		result = commands.HandleSIsMember(args)
	case "SREM":
		result = commands.HandleSRem(args)
	case "SPOP":
		result = commands.HandleSPop(args)
	case "SMOVE":
		result = commands.HandleSMove(args)
	case "SDIFF":
		result = commands.HandleSDiff(args)
	case "SDIFFSTORE":
		result = commands.HandleSDiffStore(args)
	case "SUNION":
		result = commands.HandleSUnion(args)
	case "SUNIONSTORE":
		result = commands.HandleSUnionStore(args)
	case "SRANDMEMBER":
		result = commands.HandleSRandMember(args)

	// Hash Commands
	case "HSET":
		result = commands.HandleHSet(args)
	case "HGET":
		result = commands.HandleHGet(args)
	case "HDEL":
		result = commands.HandleHDel(args)
	case "HEXISTS":
		result = commands.HandleHExists(args)
	case "HLEN":
		result = commands.HandleHLen(args)

	// Server Management Commands
	case "PING":
		result = commands.HandlePing(args)
	case "INFO":
		result = commands.HandleInfo(args)
	case "TIME":
		result = commands.HandleTime(args)
	case "DBSIZE":
		result = commands.HandleDBSize(args)
	case "FLUSHALL":
		result = commands.HandleFlushAll(args)
	case "BGSAVE":
		result = commands.HandleBGSave(args)
	case "BGREWRITEAOF":
		result = commands.HandleBGRewriteAOF(args)

	default:
		result = protocol.ErrorValue(fmt.Sprintf("Unknown command: %s", commandName))
		success = false
	}

	// Check if result is an error
	if _, isError := result.(protocol.ErrorValue); isError {
		success = false
	}

	// Convert result to string
	resultStr := s.formatCommandResult(result)

	// Broadcast command execution to WebSocket clients
	s.broadcastCommandExecution(req.Command, resultStr, success)

	cmdResult := CommandResult{
		Command:   req.Command,
		Result:    resultStr,
		Timestamp: time.Now().Unix(),
		Success:   success,
	}

	s.writeJSON(w, APIResponse{
		Success: true,
		Data:    cmdResult,
	})
}

func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, map[string]interface{}{
		"status":  "healthy",
		"version": "v0.1.0",
		"uptime":  time.Since(startTime).Seconds(),
	})
}
