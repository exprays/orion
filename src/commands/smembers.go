package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSMembers handles the SMEMBERS command
func HandleSMembers(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 1 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'smembers' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	members := data.Store.SMembers(string(key))

	response := make(protocol.ArrayValue, len(members))
	for i, member := range members {
		response[i] = protocol.BulkStringValue(member)
	}

	return response
}
