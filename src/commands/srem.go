package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSRem handles the SREM command
func HandleSRem(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'srem' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	members := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		member, ok := arg.(protocol.BulkStringValue)
		if !ok {
			return protocol.ErrorValue("ERR invalid member")
		}
		members[i] = string(member)
	}

	removed := data.Store.SRem(string(key), members...)

	return protocol.IntegerValue(removed)
}
