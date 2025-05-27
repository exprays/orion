package commands

import (
	"orion/src/aof"
	"orion/src/data"
	"orion/src/protocol"
)

// HandleHDel deletes one or more fields from a hash stored at key
func HandleHDel(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'hdel' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	// Extract field names
	fields := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		field, ok := arg.(protocol.BulkStringValue)
		if !ok {
			return protocol.ErrorValue("ERR invalid field")
		}
		fields[i] = string(field)
	}

	deleted := data.Store.HDel(string(key), fields...)

	// Log to AOF if any fields were deleted
	if deleted > 0 {
		if err := aof.AppendCommand(protocol.ArrayValue(args)); err != nil {
			return protocol.ErrorValue("ERR failed to write to AOF")
		}
	}

	return protocol.IntegerValue(deleted)
}
