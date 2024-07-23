package server

import (
	"fmt"
	"orion/src/commands"
	"orion/src/protocol"
	"strings"
)

// CommandHandler is the function signature for command handlers
type CommandHandler func(args []protocol.ORSPValue) protocol.ORSPValue

// CommandMap maps command names to their handlers
var CommandMap = map[string]CommandHandler{
	"PING":        commands.HandlePing,
	"SET":         commands.HandleSet,
	"GET":         commands.HandleGet,
	"FLUSHALL":    commands.HandleFlushAll,
	"APPEND":      commands.HandleAppend,
	"GETDEL":      commands.HandleGetDel,
	"GETEX":       commands.HandleGetEx,
	"GETSET":      commands.HandleGetSet,
	"GETRANGE":    commands.HandleGetRange,
	"INCR":        commands.HandleIncr,
	"INCRBY":      commands.HandleIncrBy,
	"INCRBYFLOAT": commands.HandleIncrByFloat,
	"LCS":         commands.HandleLCS,
	"INFO":        commands.HandleInfo,
	"TIME":        commands.HandleTime,
	"TTL":         commands.HandleTTL,

	//set commands
	"SADD":  commands.HandleSAdd,
	"SCARD": commands.HandleSCard,
}

// HandleCommand routes the command to the correct handler
func HandleCommand(command protocol.ArrayValue) protocol.ORSPValue {
	if len(command) == 0 {
		return protocol.ErrorValue("Empty command")
	}

	cmdVal, ok := command[0].(protocol.BulkStringValue)
	if !ok {
		return protocol.ErrorValue("Invalid command format")
	}

	cmd := strings.ToUpper(string(cmdVal))
	handler, exists := CommandMap[cmd]
	if !exists {
		return protocol.ErrorValue(fmt.Sprintf("Unknown command: %s", cmd))
	}

	return handler(command[1:])
}
