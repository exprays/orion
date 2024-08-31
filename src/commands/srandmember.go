package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strconv"
)

// HandleSRandMember handles the SRANDMEMBER command
func HandleSRandMember(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 1 || len(args) > 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'srandmember' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	count := 1
	if len(args) == 2 {
		countArg, ok := args[1].(protocol.BulkStringValue)
		if !ok {
			return protocol.ErrorValue("ERR invalid count")
		}
		var err error
		count, err = strconv.Atoi(string(countArg))
		if err != nil {
			return protocol.ErrorValue("ERR invalid count")
		}
	}

	result := data.Store.SRandMember(string(key), count)
	if len(result) == 0 {
		return protocol.NullValue{}
	}

	if len(args) == 1 {
		return protocol.BulkStringValue(result[0])
	}

	response := make(protocol.ArrayValue, len(result))
	for i, member := range result {
		response[i] = protocol.BulkStringValue(member)
	}
	return response
}
