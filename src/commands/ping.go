// commands/ping.go - PING command handler

package commands

import (
	"orion/src/protocol"
)

// HandlePing responds with "PONG"
func HandlePing(args []protocol.ORSPValue) protocol.ORSPValue {
	if len(args) == 0 {
		return protocol.SimpleStringValue("PONG")
	}
	return args[0]
}
