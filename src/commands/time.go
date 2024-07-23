package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strconv"
	"strings"
)

// HandleTime retrieves the current server time using ORSP
func HandleTime(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 0 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'time' command")
	}

	response, err := data.Store.Time()
	if err != nil {
		return protocol.ErrorValue("ERR " + err.Error())
	}

	// Parse the response, which is in the format "[seconds microseconds]"
	response = strings.Trim(response, "[]")
	parts := strings.Split(response, " ")
	if len(parts) != 2 {
		return protocol.ErrorValue("ERR invalid time format")
	}

	seconds, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return protocol.ErrorValue("ERR invalid seconds value")
	}

	microseconds, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return protocol.ErrorValue("ERR invalid microseconds value")
	}

	return protocol.ArrayValue{
		protocol.IntegerValue(seconds),
		protocol.IntegerValue(microseconds),
	}
}
