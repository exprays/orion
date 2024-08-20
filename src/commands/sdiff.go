package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSDiff handles the SDIFF command
func HandleSDiff(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'sdiff' command")
	}

	keys := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		key, ok := arg.(protocol.BulkStringValue)
		if !ok {
			return protocol.ErrorValue("ERR invalid key")
		}
		keys[i] = string(key)
	}

	result := data.Store.SDiff(keys...)

	response := make(protocol.ArrayValue, len(result))
	for i, member := range result {
		response[i] = protocol.BulkStringValue(member)
	}

	return response
}
