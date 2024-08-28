package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSIsMember handles the SISMEMBER command
func HandleSIsMember(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'sismember' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	member, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid member")
	}

	isMember := data.Store.SIsMember(string(key), string(member))

	if isMember {
		return protocol.IntegerValue(1)
	}
	return protocol.IntegerValue(0)
}
