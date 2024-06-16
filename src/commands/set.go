package commands

import (
	"orion/src/data"
	"strings"
)

// HandleSet sets a key-value pair in the data store and returns a RESP Simple String
func HandleSet(args []string) string {
	if len(args) != 3 {
		return "-ERROR Usage: SET key value\r\n"
	}
	key := args[1]
	value := args[2]

	// Handle the case where the value is quoted
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) > 1 {
		value = value[1 : len(value)-1]
	}

	data.Store.Set(key, value)
	return "+OK\r\n"
}
