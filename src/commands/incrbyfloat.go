// commands/incrbyfloat.go

package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strconv"
)

// HandleIncrByFloat increments the float value of a key by a specified amount
func HandleIncrByFloat(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'incrbyfloat' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	incrementStr, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid increment")
	}

	increment, err := strconv.ParseFloat(string(incrementStr), 64)
	if err != nil {
		return protocol.ErrorValue("ERR increment is not a valid float")
	}

	newValue, err := data.Store.IncrByFloat(string(key), increment)
	if err != nil {
		return protocol.ErrorValue("ERR " + err.Error())
	}

	return protocol.BulkStringValue(strconv.FormatFloat(newValue, 'f', -1, 64))
}
