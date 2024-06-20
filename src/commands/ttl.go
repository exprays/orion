// commands/ttl.go - TTL command handler

package commands

import (
	"orion/src/data"
	"strconv"
)

// HandleTTL retrieves the TTL for a given key
func HandleTTL(args []string) string {
	if len(args) != 1 {
		return "ERROR: Usage: TTL key"
	}
	key := args[0]
	ttl := data.Store.TTL(key)
	if ttl == -1 {
		return "ERROR: Key does not exist or has no associated TTL"
	}
	return strconv.FormatInt(ttl, 10)
}
