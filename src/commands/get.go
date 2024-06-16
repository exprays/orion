// commands/get.go - GET command handler

package commands

import (
	"orion/src/data"
	"strconv"
)

// HandleGet retrieves the value for a key from the data store and returns a RESP Bulk String
func HandleGet(args []string) string {
	if len(args) != 1 {
		return "-ERROR Usage: GET key\r\n"
	}
	key := args[0]
	value, ok := data.Store.Get(key)
	if !ok {
		return "$-1\r\n" // Return a Null Bulk String ("-1\r\n")
	}
	return "$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n"
}
