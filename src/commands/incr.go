// commands/incr.go

package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleIncr increments the integer value of a key by 1
func HandleIncr(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 1 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'incr' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	newValue, err := data.Store.Incr(string(key))
	if err != nil {
		return protocol.ErrorValue("ERR " + err.Error())
	}

	return protocol.IntegerValue(newValue)
}
