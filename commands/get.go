package commands

import (
    "orion/data"
)

// HandleGet retrieves the value for a key from the data store
func HandleGet(args []string) string {
    if len(args) != 1 {
        return "ERROR: Usage: GET key"
    }
    key := args[0]
    value, ok := data.Store.Get(key) // Use data.Store to access the global store instance
    if !ok {
        return "nil"
    }
    return value
}
