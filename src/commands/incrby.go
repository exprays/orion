// commands/incrby.go - INCRBY command handler

package commands

import (
	"fmt"
	"orion/src/data"
	"strconv"
)

// HandleIncrBy increments the integer value of a key by a specified amount
func HandleIncrBy(args []string) string {
	if len(args) != 2 {
		return "ERROR: Usage: INCRBY key increment"
	}
	key := args[0]
	increment, err := strconv.Atoi(args[1])
	if err != nil {
		return "ERROR: Increment must be an integer"
	}

	newValue, err := data.Store.IncrBy(key, increment)
	if err != nil {
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	return strconv.Itoa(newValue)
}
