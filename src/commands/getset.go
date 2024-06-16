// commands/getset.go - Implement the GETSET command

package commands

import (
	"orion/src/data"
)

// HandleGetSet sets a new value for a key and returns its old value
func HandleGetSet(args []string) string {
	if len(args) != 2 {
		return "ERROR: Usage: GETSET key value"
	}
	key := args[0]
	value := args[1]
	oldValue, exists := data.Store.GetSet(key, value)
	if !exists {
		return "(nil)"
	}
	return oldValue
}
