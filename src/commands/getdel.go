// commands/getdel.go

package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleGetDel retrieves the value of a key and deletes it from the data store
func HandleGetDel(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 1 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'getdel' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	value, exists := data.Store.GetDel(string(key))
	if !exists {
		return protocol.NullValue{}
	}

	return protocol.BulkStringValue(value)
}
