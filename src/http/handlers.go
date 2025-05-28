package http

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *HTTPServer) handleCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Invalid request body",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// For now, return a simple response
	result := CommandResult{
		Command:   req.Command,
		Result:    "Command execution not implemented yet",
		Timestamp: time.Now().Unix(),
		Success:   false,
	}

	response := APIResponse{
		Success: false,
		Data:    result,
		Error:   "Command execution not implemented",
	}

	json.NewEncoder(w).Encode(response)
}

func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(startTime).String(),
	}

	json.NewEncoder(w).Encode(response)
}
