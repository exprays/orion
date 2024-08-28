package commands

import (
	"orion/src/data"
	"orion/src/protocol"
	"strings"
)

// HandleInfo returns server statistics in a formatted string using ORSP
func HandleInfo(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) != 0 {
		return protocol.ErrorValue("ERR wrong number of arguments for 'info' command")
	}

	info := data.Store.Info()
	sections := strings.Split(info, "\n")

	response := make(protocol.ArrayValue, 0, len(sections))
	for _, section := range sections {
		if strings.TrimSpace(section) != "" {
			response = append(response, protocol.BulkStringValue(section))
		}
	}

	return response
}
