package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleHLen returns the number of fields in a hash stored at key
func HandleHLen(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 1 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'hlen' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	length := data.Store.HLen(string(key))

	return protocol.IntegerValue(length)
}
