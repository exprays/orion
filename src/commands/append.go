package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleAppend appends a value to an existing string key
func HandleAppend(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'append' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	value, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid value")
	}

	// Append the value to the existing value in the store
	data.Store.Append(string(key), string(value))

	// Get the new length of the string
	newValue, _ := data.Store.Get(string(key))
	length := len(newValue)

	// Return the length as an IntegerValue
	return protocol.IntegerValue(length)
}
