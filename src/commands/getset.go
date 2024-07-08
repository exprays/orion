// commands/getset.go

package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleGetSet sets a new value for a key and returns its old value
func HandleGetSet(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'getset' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	newValue, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid value")
	}

	oldValue, exists := data.Store.GetSet(string(key), string(newValue))
	if !exists {
		return protocol.NullValue{}
	}

	return protocol.BulkStringValue(oldValue)
}
