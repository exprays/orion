package hunter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	maxHistorySize = 1000
	historyDir     = ".orion_history"
	historyFile    = "command_history.json"
)

// HistoryEntry represents a single command in the history
type HistoryEntry struct {
	Command   string    `json:"command"`
	Timestamp time.Time `json:"timestamp"`
}

// CommandHistory manages the command history
type CommandHistory struct {
	Entries  []HistoryEntry `json:"entries"`
	Position int            `json:"-"` // Current position when navigating (not persisted)
}

// NewCommandHistory initializes a new command history
func NewCommandHistory() *CommandHistory {
	history := &CommandHistory{
		Entries:  []HistoryEntry{},
		Position: -1,
	}
	history.Load()
	return history
}

// Add adds a command to the history
func (h *CommandHistory) Add(command string) {
	// Don't add empty commands or duplicates of the last command
	if command == "" || (len(h.Entries) > 0 && h.Entries[len(h.Entries)-1].Command == command) {
		return
	}

	// Add the command
	h.Entries = append(h.Entries, HistoryEntry{
		Command:   command,
		Timestamp: time.Now(),
	})

	// Trim if we exceed max size
	if len(h.Entries) > maxHistorySize {
		h.Entries = h.Entries[len(h.Entries)-maxHistorySize:]
	}

	// Reset position to the end
	h.Position = len(h.Entries)

	// Save to file
	h.Save()
}

// Previous returns the previous command in history
func (h *CommandHistory) Previous() (string, bool) {
	if len(h.Entries) == 0 {
		return "", false
	}

	if h.Position > 0 {
		h.Position--
		return h.Entries[h.Position].Command, true
	} else if h.Position == 0 {
		// Already at the beginning
		return h.Entries[0].Command, true
	}

	// Initialize position to the end if it's -1
	if h.Position == -1 {
		h.Position = len(h.Entries) - 1
		return h.Entries[h.Position].Command, true
	}

	return "", false
}

// Next returns the next command in history
func (h *CommandHistory) Next() (string, bool) {
	if len(h.Entries) == 0 || h.Position == -1 {
		return "", false
	}

	if h.Position < len(h.Entries)-1 {
		h.Position++
		return h.Entries[h.Position].Command, true
	}

	// At the end, return empty and reset position
	h.Position = len(h.Entries)
	return "", true
}

// List returns the last n commands in history
func (h *CommandHistory) List(n int) []HistoryEntry {
	if n <= 0 || len(h.Entries) == 0 {
		return []HistoryEntry{}
	}

	if n > len(h.Entries) {
		n = len(h.Entries)
	}

	return h.Entries[len(h.Entries)-n:]
}

// Save persists the command history to disk
func (h *CommandHistory) Save() error {
	// Create directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	historyDirPath := filepath.Join(homeDir, historyDir)
	if err := os.MkdirAll(historyDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %v", err)
	}

	// Marshal to JSON
	data, err := json.Marshal(h)
	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}

	// Write to file
	historyFilePath := filepath.Join(historyDirPath, historyFile)
	if err := os.WriteFile(historyFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %v", err)
	}

	return nil
}

// Load reads the command history from disk
func (h *CommandHistory) Load() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	historyFilePath := filepath.Join(homeDir, historyDir, historyFile)
	data, err := os.ReadFile(historyFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist yet, that's okay
			return nil
		}
		return fmt.Errorf("failed to read history file: %v", err)
	}

	if err := json.Unmarshal(data, h); err != nil {
		return fmt.Errorf("failed to unmarshal history: %v", err)
	}

	// Initialize position to the end
	h.Position = len(h.Entries)
	return nil
}
