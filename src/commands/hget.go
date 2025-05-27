package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleHGet gets the value of a field from a hash stored at key
func HandleHGet(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'hget' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	field, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid field")
	}

	value, exists := data.Store.HGet(string(key), string(field))
	if !exists {
		return protocol.NullValue{}
	}

	return protocol.BulkStringValue(value)
}
