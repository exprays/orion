package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSCard returns the cardinality (number of elements) of the set stored at key
func HandleSCard(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 1 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'scard' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	cardinality := data.Store.SCard(string(key))
	return protocol.IntegerValue(cardinality)
}
