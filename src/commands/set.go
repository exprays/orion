package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strconv"
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
			return protocol.ErrorValue("ERR syntax error")
		}
		switch string(arg) {
		case "EX":
			if i+1 < len(args) {
				if seconds, ok := args[i+1].(protocol.BulkStringValue); ok {
					s, err := strconv.Atoi(string(seconds))
					if err != nil {
						return protocol.ErrorValue("ERR value is not an integer or out of range")
					}
					expiration = time.Duration(s) * time.Second
					i++
				} else {
					return protocol.ErrorValue("ERR syntax error")
				}
			} else {
				return protocol.ErrorValue("ERR syntax error")
			}
		case "PX":
			if i+1 < len(args) {
				if milliseconds, ok := args[i+1].(protocol.BulkStringValue); ok {
					ms, err := strconv.Atoi(string(milliseconds))
					if err != nil {
						return protocol.ErrorValue("ERR value is not an integer or out of range")
					}
					expiration = time.Duration(ms) * time.Millisecond
					i++
				} else {
					return protocol.ErrorValue("ERR syntax error")
				}
			} else {
				return protocol.ErrorValue("ERR syntax error")
			}
		case "XX":
			xx = true
		case "NX":
			nx = true
		default:
			return protocol.ErrorValue("ERR syntax error")
		}
	}

	// Check XX and NX conditions
	exists := data.Store.Exists(string(key))
	if (xx && !exists) || (nx && exists) {
		return protocol.NullValue{}
	}

	// Set the value
	data.Store.Set(string(key), string(value), expiration)

	return protocol.SimpleStringValue("OK")
}
