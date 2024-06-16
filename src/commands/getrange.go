// commands/getrange.go - Implement the GETRANGE command

package commands

import (
	"orion/src/data"
	"strconv"
)

// HandleGetRange retrieves a substring of the string value stored at a key
func HandleGetRange(args []string) string {
	if len(args) != 3 {
		return "ERROR: Usage: GETRANGE key start end"
	}
	key := args[0]
	start, err1 := strconv.Atoi(args[1])
	end, err2 := strconv.Atoi(args[2])
	if err1 != nil || err2 != nil {
		return "ERROR: Invalid start or end index"
	}
	value := data.Store.GetRange(key, start, end)
	return value
}
