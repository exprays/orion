package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"time"
)

// HandleSet sets a key-value pair in the data store
func HandleSet(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) < 2 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'set' command")
	}

	key, ok := args[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid key")
	}

	value, ok := args[1].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("ERR invalid value")
	}

	// Parse optional arguments
	var expiration time.Duration
	var xx, nx bool
	for i := 2; i < len(args); i++ {
		arg, ok := args[i].(protocol.BulkStringValue)
		if !ok {
			continue
		}
		switch string(arg) {
		case "EX":
			if i+1 < len(args) {
				if seconds, ok := args[i+1].(protocol.IntegerValue); ok {
					expiration = time.Duration(seconds) * time.Second
					i++
				}
			}
		case "PX":
			if i+1 < len(args) {
				if milliseconds, ok := args[i+1].(protocol.IntegerValue); ok {
					expiration = time.Duration(milliseconds) * time.Millisecond
					i++
				}
			}
		case "XX":
			xx = true
		case "NX":
			nx = true
		}
	}

	// Check XX and NX conditions
	exists := data.Store.Exists(string(key))
	if (xx && !exists) || (nx && exists) {
		return protocol.NullValue{}
	}

	// Set the value using the updated Set method
	data.Store.Set(string(key), string(value), expiration)

	return protocol.SimpleStringValue("OK")
}
