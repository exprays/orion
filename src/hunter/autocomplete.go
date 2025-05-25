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

	// Client commands
	"CLEAR", "HISTORY", "HELP", "EXIT", "QUIT",
}

// newAutoCompleter creates a new completer for readline
func newAutoCompleter() *readline.PrefixCompleter {
	var items []readline.PrefixCompleterInterface

	for _, cmd := range commandList {
		items = append(items, readline.PcItem(cmd))
	}

	return readline.NewPrefixCompleter(items...)
}

// CompleteCommand provides completion for command names
func CompleteCommand(line string, pos int) (string, []string, string) {
	line = strings.ToUpper(line)
	parts := strings.Fields(line[:pos])

	if len(parts) == 0 {
		// Complete command at the beginning
		var candidates []string
		for _, cmd := range commandList {
			candidates = append(candidates, cmd)
		}
		return "", candidates, ""
	}

	if len(parts) == 1 && !strings.HasSuffix(line[:pos], " ") {
		// Complete the first word (command)
		var candidates []string
		prefix := parts[0]
		for _, cmd := range commandList {
			if strings.HasPrefix(cmd, prefix) {
				candidates = append(candidates, cmd)
			}
		}
		return prefix, candidates, ""
	}

	return "", nil, ""
}
