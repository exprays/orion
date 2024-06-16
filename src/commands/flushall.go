// commands/flushall.go - FLUSHALL command handler

package commands

import "orion/src/data"

// HandleFlushAll clears all key-value pairs from the data store
func HandleFlushAll(args []string) string {
	if len(args) != 0 {
		return "ERROR: FLUSHALL does not accept arguments"
	}
	data.Store.FlushAll()
	return "OK"
}
