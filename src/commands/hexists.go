package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleHExists checks if a field exists in a hash stored at key
func HandleHExists(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'hexists' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	field, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid field")
	}

	exists := data.Store.HExists(string(key), string(field))

	if exists {
		return protocol.IntegerValue(1)
	}
	return protocol.IntegerValue(0)
}
