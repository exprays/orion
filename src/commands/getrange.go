// commands/getrange.go

package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strconv"
)

// HandleGetRange retrieves a substring of the string value stored at a key
func HandleGetRange(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 3 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'getrange' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	startStr, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid start index")
	}

	endStr, ok := args[2].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid end index")
	}

	start, err := strconv.Atoi(string(startStr))
	if err != nil {
		return protocol.ErrorValue("ERR invalid start index")
	}

	end, err := strconv.Atoi(string(endStr))
	if err != nil {
		return protocol.ErrorValue("ERR invalid end index")
	}

	value := data.Store.GetRange(string(key), start, end)

	return protocol.BulkStringValue(value)
}
