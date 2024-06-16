// commands/decr.go

package commands

import (
    "orion/data"
    "strconv"
)

// HandleDecr decrements the value of a string key by 1
func HandleDecr(args []string) string {
    if len(args) != 1 {
        return "ERROR: Usage: DECR key"
    }
    key := args[0]
    newValue := data.Store.DecrBy(key, 1)
    return strconv.Itoa(newValue)
}
