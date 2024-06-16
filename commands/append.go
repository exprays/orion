// commands/append.go - APPEND command handler

package commands

import (
    "orion/data"
    "strconv"
)

// HandleAppend appends a value to an existing string key
func HandleAppend(args []string) string {
    if len(args) != 2 {
        return "ERROR: Usage: APPEND key value"
    }
    key, value := args[0], args[1]

    // Append the value to the existing value in the store
    data.Store.Append(key, value)

    // Get the new length of the string
    newValue, _ := data.Store.Get(key)
    length := len(newValue)

    // Return the length as a string
    return strconv.Itoa(length)
}
