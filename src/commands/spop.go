package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strconv"
)

// HandleSPop handles the SPOP command
func HandleSPop(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 1 || len(args) > 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'spop' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	count := 1
	if len(args) == 2 {
		countStr, ok := args[1].(protocol.BulkStringValue)
		if !ok {
			return protocol.ErrorValue("ERR invalid count")
		}
		var err error
		count, err = strconv.Atoi(string(countStr))
		if err != nil || count < 0 {
			return protocol.ErrorValue("ERR invalid count")
		}
	}

	poppedMembers := data.Store.SPop(string(key), count)

	if len(poppedMembers) == 0 {
		return protocol.NullValue{}
	}

	if count == 1 {
		return protocol.BulkStringValue(poppedMembers[0])
	}

	response := make(protocol.ArrayValue, len(poppedMembers))
	for i, member := range poppedMembers {
		response[i] = protocol.BulkStringValue(member)
	}
	return response
}
