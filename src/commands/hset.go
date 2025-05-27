package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleHSet sets field-value pairs in a hash stored at key
func HandleHSet(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 3 || len(args)%2 == 0 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'hset' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	// Extract field-value pairs
	fieldValues := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		value, ok := arg.(protocol.BulkStringValue)
		if !ok {
			return protocol.ErrorValue("ERR invalid field or value")
		}
		fieldValues[i] = string(value)
	}

	created := data.Store.HSet(string(key), fieldValues...)
	if created < 0 {
		return protocol.ErrorValue("ERR invalid number of field-value pairs")
	}
	return protocol.IntegerValue(created)
}
