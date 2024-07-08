// commands/getex.go

package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strconv"
)

// HandleGetEx retrieves the value of a key and sets its expiration in seconds
func HandleGetEx(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'getex' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	secondsStr, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid seconds argument")
	}

	seconds, err := strconv.ParseInt(string(secondsStr), 10, 64)
	if err != nil {
		return protocol.ErrorValue("ERR invalid seconds argument")
	}

	value, exists := data.Store.GetEx(string(key), seconds)
	if !exists {
		return protocol.NullValue{}
	}

	return protocol.BulkStringValue(value)
}
