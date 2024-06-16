// server/handler.go
package server

import (
	"orion/src/commands"
	"strings"
)

// CommandHandler is the function signature for command handlers
type CommandHandler func(args []string) string

// CommandMap maps command names to their handlers
var CommandMap = map[string]CommandHandler{
	"PING":     commands.HandlePing,
	"SET":      commands.HandleSet,
	"GET":      commands.HandleGet,
	"FLUSHALL": commands.HandleFlushAll,
	"APPEND":   commands.HandleAppend,
	"DECR":     commands.HandleDecr,
	"DECRBY":   commands.HandleDecrBy,
	"GETDEL":   commands.HandleGetDel,
	"GETEX":    commands.HandleGetEx,
	"GETRANGE": commands.HandleGetRange,
	"GETSET":   commands.HandleGetSet,
	// Add more THUNDER type handlers as needed
}

// HandleCommand routes the command to the correct handler and formats the response in Thunder protocol
func HandleCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "-ERROR Empty command\r\n"
	}

	cmd := strings.ToUpper(parts[0])
	handler, exists := CommandMap[cmd]
	if !exists {
		return "-ERROR Unknown command\r\n"
	}

	return formatThunder(handler(parts[1:]))
}

// formatThunder formats a response as a THUNDER response
func formatThunder(response string) string {
	if strings.HasPrefix(response, "ERROR:") {
		return "-" + response[6:] + "\r\n"
	}
	return "+" + response + "\r\n"
}
