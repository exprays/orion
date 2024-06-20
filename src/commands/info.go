// commands/info.go - INFO command handler

package commands

import (
	"orion/src/data"
)

// HandleInfo returns server statistics in a formatted string
func HandleInfo(args []string) string {
	if len(args) != 0 {
		return "ERROR: Usage: INFO"
	}
	info := data.Store.Info() // Use data.Store to access the global store instance
	return info
}
