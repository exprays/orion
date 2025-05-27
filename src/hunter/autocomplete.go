package hunter

import (
	"strings"

	"github.com/chzyer/readline"
)

// commandList contains all the available commands for autocomplete
// update this list as new commands are added to the CLI
var commandList = []string{
	// Server Management commands
	"BGSAVE", "BGREWRITEAOF", "FLUSHALL", "PING", "TIME", "INFO", "DBSIZE",

	// String commands
	"SET", "GET", "APPEND", "GETDEL", "GETEX", "GETSET", "GETRANGE",
	"INCR", "INCRBY", "INCRBYFLOAT", "LCS", "TTL",

	// Set commands
	"SADD", "SCARD", "SMEMBERS", "SISMEMBER", "SREM", "SPOP", "SMOVE",
	"SDIFF", "SDIFFSTORE", "SUNION", "SUNIONSTORE", "SRANDMEMBER",

	// Hash commands
	"HSET", "HGET", "HDEL", "HEXISTS", "HLEN",

	// Client commands
	"CLEAR", "HISTORY", "HELP", "EXIT", "QUIT",
}

// newAutoCompleter creates a new completer for readline
func newAutoCompleter() *readline.PrefixCompleter {
	var items []readline.PrefixCompleterInterface

	for _, cmd := range commandList {
		// Add both uppercase and lowercase versions for case-insensitive matching
		items = append(items, readline.PcItem(cmd))
		items = append(items, readline.PcItem(strings.ToLower(cmd)))
	}

	return readline.NewPrefixCompleter(items...)
}
