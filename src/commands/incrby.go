// commands/incrby.go

package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strconv"
)

// HandleIncrBy increments the integer value of a key by a specified amount
func HandleIncrBy(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'incrby' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	incrementStr, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid increment")
	}

	increment, err := strconv.Atoi(string(incrementStr))
	if err != nil {
		return protocol.ErrorValue("ERR increment must be an integer")
	}

	newValue, err := data.Store.IncrBy(string(key), increment)
	if err != nil {
		return protocol.ErrorValue("ERR " + err.Error())
	}

	return protocol.IntegerValue(newValue)
}
