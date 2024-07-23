package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSAdd adds the specified members to the set stored at key
func HandleSAdd(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'sadd' command")
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

	added := data.Store.SAdd(string(key), members...)
	return protocol.IntegerValue(added)
}
