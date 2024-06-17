package commands

import (
	"orion/src/data"
	"strings"
)

// HandleSet sets a key-value pair in the data store
func HandleSet(args []string) string {
	if len(args) < 2 {
		return "ERROR: Usage: SET key value"
	}

	key := args[0]
	value := strings.Join(args[1:], " ") // Join all remaining arguments into a single value
	data.Store.Set(key, value)
	return "OK"
}
