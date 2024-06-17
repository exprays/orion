// commands/getset.go - Implement the GETSET command

package commands

import (
	"orion/src/data"
	"strings"
)

// HandleGetSet sets a new value for a key and returns its old value
func HandleGetSet(args []string) string {
	if len(args) < 2 {
		return "ERROR: Usage: GETSET key value"
	}
	key := args[0]
	value := strings.Join(args[1:], " ") // Join all remaining arguments into a single value
	oldValue, exists := data.Store.GetSet(key, value)
	if !exists {
		return "(nil)"
	}
	return oldValue
}
