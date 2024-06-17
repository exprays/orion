// commands/incr.go - INCR command handler

package commands

import (
	"fmt"
	"orion/src/data"
	"strconv"
)

// HandleIncr increments the integer value of a key by 1
func HandleIncr(args []string) string {
	if len(args) != 1 {
		return "ERROR: Usage: INCR key"
	}
	key := args[0]

	newValue, err := data.Store.Incr(key)
	if err != nil {
		return fmt.Sprintf("ERROR: %s", err.Error())
	}
	return strconv.Itoa(newValue)
}
