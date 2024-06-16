// commands/set.go - SET command handler

package commands

import (
	"orion/src/data"
)

// HandleSet sets a key-value pair in the data store
func HandleSet(args []string) string {
	if len(args) != 2 {
		return "ERROR: Usage: SET key value"
	}
	key, value := args[0], args[1]
	data.Store.Set(key, value) // Use data.Store to access the global store instance
	return "OK"
}
