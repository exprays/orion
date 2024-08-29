package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSDiffStore handles the SDIFFSTORE command
func HandleSDiffStore(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 3 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'sdiffstore' command")
	}

	destination, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid destination key")
	}

	keys := make([]string, len(args)-2)
	for i, arg := range args[2:] {
		key, ok := arg.(protocol.BulkStringValue)
		if !ok {
			return protocol.ErrorValue("ERR invalid key")
		}
		keys[i] = string(key)
	}

	count := data.Store.SDiffStore(string(destination), keys...)

	return protocol.IntegerValue(int64(count))
}
