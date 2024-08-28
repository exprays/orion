package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleTTL retrieves the TTL for a given key using ORSP
func HandleTTL(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 1 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'ttl' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	ttl := data.Store.TTL(string(key))

	switch {
	case ttl == -2:
		// Key does not exist
		return protocol.IntegerValue(-2)
	case ttl == -1:
		// Key exists but has no associated TTL
		return protocol.IntegerValue(-1)
	default:
		// Key exists and has a TTL
		return protocol.IntegerValue(ttl)
	}
}
