// commands/getex.go - Implement the GETEX command

package commands

import (
    "orion/data"
    "strconv"
)

// HandleGetEx retrieves the value of a key and sets expiration in seconds
func HandleGetEx(args []string) string {
    if len(args) != 2 {
        return "ERROR: Usage: GETEX key seconds"
    }
    key := args[0]
    seconds, err := strconv.ParseInt(args[1], 10, 64)
    if err != nil {
        return "ERROR: Invalid seconds argument"
    }
    value, exists := data.Store.GetEx(key, seconds)
    if !exists {
        return "(nil)"
    }
    return value
}
