package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleSMove handles the SMOVE command
func HandleSMove(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 3 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'smove' command")
	}

	source, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid source key")
	}

	destination, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid destination key")
	}

	member, ok := args[2].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid member")
	}

	moved := data.Store.SMove(string(source), string(destination), string(member))

	if moved {
		return protocol.IntegerValue(1)
	}
	return protocol.IntegerValue(0)
}
