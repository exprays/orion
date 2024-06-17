// commands/incrbyfloat.go - INCRBYFLOAT command handler

package commands

import (
	"orion/src/data"
	"strconv"
)

// HandleIncrByFloat increments the float value of a key by a specified amount
func HandleIncrByFloat(args []string) string {
	if len(args) != 2 {
		return "ERROR: Usage: INCRBYFLOAT key increment"
	}
	key := args[0]
	increment, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return "ERROR: Increment must be a float"
	}

	newValue, err := data.Store.IncrByFloat(key, increment)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return strconv.FormatFloat(newValue, 'f', -1, 64)
}
