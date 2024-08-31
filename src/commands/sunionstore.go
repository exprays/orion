package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSUnionStore handles the SUNIONSTORE command
func HandleSUnionStore(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'sunionstore' command")
	}

	destination, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid destination key")
	}

	keys := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		key, ok := arg.(protocol.BulkStringValue)
		if !ok {
			return protocol.ErrorValue("ERR invalid key")
		}
		keys[i] = string(key)
	}

	count := data.Store.SUnionStore(string(destination), keys...)
	return protocol.IntegerValue(count)
}
