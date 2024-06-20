// commands/setex.go - SETEX command handler

package commands

import (
	"orion/src/data"
	"strconv"
)

// HandleSetEx sets a key to a value and sets a TTL in seconds
func HandleSetEx(args []string) string {
	if len(args) != 3 {
		return "ERROR: Usage: SETEX key seconds value"
	}
	key := args[0]
	seconds, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || seconds <= 0 {
		return "ERROR: Seconds must be a positive integer"
	}
	value := args[2]

	data.Store.SetEx(key, value, seconds)
	return "OK"
}
