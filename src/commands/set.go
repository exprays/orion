// commands/set.go
package commands

import (
	"orion/src/data"
)

// HandleSet sets a key-value pair in the data store and returns a Thunder Simple String
func HandleSet(args []string) string {
	if len(args) != 2 {
		return "-ERROR Usage: SET key value\r\n"
	}
	key, value := args[0], args[1]
	data.Store.Set(key, value)
	return "OK"
}
