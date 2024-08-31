package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSUnion handles the SUNION command
func HandleSUnion(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 1 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'sunion' command")
	}

	keys := make([]string, len(args))
	for i, arg := range args {
		key, ok := arg.(protocol.BulkStringValue)
		if !ok {
			return protocol.ErrorValue("ERR invalid key")
		}
		keys[i] = string(key)
	}

	result := data.Store.SUnion(keys...)
	response := make(protocol.ArrayValue, len(result))
	for i, member := range result {
		response[i] = protocol.BulkStringValue(member)
	}
	return response
}
