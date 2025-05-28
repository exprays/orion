package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"orion/src/core"
	"orion/src/protocol"
)

func (s *HTTPServer) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *HTTPServer) writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: false,
		Error:   message,
	}

	json.NewEncoder(w).Encode(response)
}

func (s *HTTPServer) getKeyInfo(key string) *KeyInfo {
	// Get key type
	keyType := core.GetKeyType(key)
	if keyType == "" {
		return nil
	}

	// Get TTL
	ttl := core.GetTTL(key)

	keyInfo := &KeyInfo{
		Key:  key,
		Type: keyType,
		TTL:  ttl,
		Size: 0, // Will be calculated based on content
	}

	switch keyType {
	case "string":
		if value := core.GetString(key); value != nil {
			keyInfo.Value = *value
			keyInfo.Size = len(*value)
		}
	case "set":
		if members := core.GetSetMembers(key); members != nil {
			keyInfo.Members = members
			keyInfo.Size = calculateSetSize(members)
		}
	case "hash":
		if fields := core.GetHashFields(key); fields != nil {
			keyInfo.Fields = fields
			keyInfo.Size = calculateHashSize(fields)
		}
	}

	return keyInfo
}

func (s *HTTPServer) formatCommandResult(result protocol.ORSPValue) string {
	switch v := result.(type) {
	case protocol.SimpleStringValue:
		return string(v)
	case protocol.ErrorValue:
		return fmt.Sprintf("ERROR: %s", string(v))
	case protocol.IntegerValue:
		return strconv.FormatInt(int64(v), 10)
	case protocol.BulkStringValue:
		return fmt.Sprintf("\"%s\"", string(v))
	case protocol.ArrayValue:
		var parts []string
		for i, item := range v {
			parts = append(parts, fmt.Sprintf("%d) %s", i+1, s.formatCommandResult(item)))
		}
		return strings.Join(parts, "\n")
	case protocol.NullValue:
		return "(nil)"
	case protocol.BooleanValue:
		if bool(v) {
			return "(integer) 1"
		}
		return "(integer) 0"
	default:
		return fmt.Sprintf("%v", result)
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func calculateSetSize(members []string) int {
	size := 0
	for _, member := range members {
		size += len(member) + 8 // 8 bytes overhead per member
	}
	return size
}

func calculateHashSize(fields map[string]string) int {
	size := 0
	for key, value := range fields {
		size += len(key) + len(value) + 16 // 16 bytes overhead per field
	}
	return size
}
