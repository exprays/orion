// commands/decrby.go - DECRBY command handler

package commands

import (
	"orion/src/data"
	"strconv"
)

// HandleDecrBy decrements the value of a string key by the given number
func HandleDecrBy(args []string) string {
	if len(args) != 2 {
		return "ERROR: Usage: DECRBY key decrement"
	}
	key := args[0]
	decrement, err := strconv.Atoi(args[1])
	if err != nil {
		return "ERROR: decrement must be an integer"
	}
	newValue := data.Store.DecrBy(key, decrement)
	return strconv.Itoa(newValue)
}
