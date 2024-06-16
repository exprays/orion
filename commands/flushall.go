// commands/flushall.go
package commands

import "orion/data"

// HandleFlushAll clears all key-value pairs from the data store
func HandleFlushAll(args []string) string {
    if len(args) != 0 {
        return "ERROR: FLUSHALL does not accept arguments"
    }
    data.Store.FlushAll()
    return "OK"
}