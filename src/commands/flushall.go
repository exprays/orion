package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleFlushAll clears all key-value pairs from the data store
func HandleFlushAll(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 0 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'flushall' command")
	}

	data.Store.FlushAll()

	return protocol.SimpleStringValue("OK")
}
