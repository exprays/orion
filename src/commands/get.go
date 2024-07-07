package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleGet retrieves the value for a key from the data store
func HandleGet(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 1 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'get' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	value, exists := data.Store.Get(string(key))
	if !exists {
		return protocol.NullValue{}
	}

	return protocol.BulkStringValue(value)
}
