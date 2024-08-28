// commands/dbsize.go

package commands

import (
	"orion/src/data"
	"orion/src/protocol"
)

// HandleDBSize returns the number of keys in the data store
func HandleDBSize(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 0 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'dbsize' command")
	}

	// Get the number of keys in the data store
	dbSize := data.Store.DBSize()
	return protocol.IntegerValue(dbSize)
}
