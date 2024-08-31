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

	//Server Management commands
	"BGSAVE":       commands.HandleBGSave,
	"BGREWRITEAOF": commands.HandleBGRewriteAOF,
	"FLUSHALL":     commands.HandleFlushAll,
	"PING":         commands.HandlePing,
	"TIME":         commands.HandleTime,
	"INFO":         commands.HandleInfo,
	"DBSIZE":       commands.HandleDBSize,

	//String commands

	"SET":         commands.HandleSet,
	"GET":         commands.HandleGet,
	"APPEND":      commands.HandleAppend,
	"GETDEL":      commands.HandleGetDel,
	"GETEX":       commands.HandleGetEx,
	"GETSET":      commands.HandleGetSet,
	"GETRANGE":    commands.HandleGetRange,
	"INCR":        commands.HandleIncr,
	"INCRBY":      commands.HandleIncrBy,
	"INCRBYFLOAT": commands.HandleIncrByFloat,
	"LCS":         commands.HandleLCS,
	"TTL":         commands.HandleTTL,

	//set commands
	"SADD":        commands.HandleSAdd,
	"SCARD":       commands.HandleSCard,
	"SMEMBERS":    commands.HandleSMembers,
	"SISMEMBER":   commands.HandleSIsMember,
	"SREM":        commands.HandleSRem,
	"SPOP":        commands.HandleSPop,
	"SMOVE":       commands.HandleSMove,
	"SDIFF":       commands.HandleSDiff,
	"SDIFFSTORE":  commands.HandleSDiffStore,
	"SUNION":      commands.HandleSUnion,
	"SUNIONSTORE": commands.HandleSUnionStore,
	"SRANDMEMBER": commands.HandleSRandMember,
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
